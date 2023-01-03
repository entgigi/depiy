package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	utility "github.com/entgigi/depiy/common-libs/utilities"
	"github.com/entgigi/depiy/operators/bundle-operator/api/v1alpha1"
	"github.com/entgigi/depiy/services/bundle-rest-app/controllers"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes/scheme"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var port string
	found := true
	if port, found = os.LookupEnv("PORT"); !found {
		port = ":8080"
	}

	if _, err := utility.GetWatchNamespace(); err != nil {
		log.Fatalf("error retrive kubernetes namespace: %s\n", err)

	}

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"data": "Hello World"})
	})

	router.GET("/readyz", controllers.Readyz)
	router.GET("/healthz", controllers.Healthz)
	router.GET("/version", controllers.GetVersion)

	v1alpha1.AddToScheme(scheme.Scheme)
	config, err := getKubeClient()
	if err != nil {
		log.Fatalf("error retrive kubernetes configuration: %s\n", err)
	}

	bundleCtrl, err := controllers.NewBundleCtrl(config)
	if err != nil {
		log.Fatalf("error create bundle ctrl : %s\n", err)
	}

	router.GET("/bundles", bundleCtrl.ListBundles)
	router.GET("/bundles/:code", bundleCtrl.GetBundle)

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}
	// https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

func getKubeClient() (*rest.Config, error) {
	var config *rest.Config
	var err error
	config, err = rest.InClusterConfig()
	if err != nil {
		if err == rest.ErrNotInCluster {
			var kubeconfig string
			if home := homedir.HomeDir(); home != "" {
				kubeconfig = filepath.Join(home, ".kube", "config")
			}

			var internalError error
			config, internalError = clientcmd.BuildConfigFromFlags("", kubeconfig)
			if internalError != nil {
				return nil, err
			}
			fmt.Println("Use kube config")
		}
	} else {
		fmt.Println("Use incluster config")
	}

	return config, nil
}
