package controllers

import (
	"fmt"
	"net/http"

	utility "github.com/entgigi/depiy/common-libs/utilities"
	"github.com/entgigi/depiy/operators/bundle-operator/api/v1alpha1"
	"github.com/entgigi/depiy/services/bundle-rest-app/clients"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type BundleDto struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Title         string `json:"title"`
	Icon          string `json:"icon"`
	SignatureInfo string `json:"signatureInfo"`
}

type BundleCtrl struct {
	apiClient *clients.BundleV1Alpha1Client
}

func NewBundleCtrl(config *rest.Config) (*BundleCtrl, error) {
	s, err := clients.NewForConfig(config)
	return &BundleCtrl{s}, err
}

func (bc *BundleCtrl) ListBundles(ctx *gin.Context) {

	ns, _ := utility.GetWatchNamespace()
	bundles, err := bc.apiClient.Bundles(ns).List(metav1.ListOptions{})
	fmt.Println("ListBundles")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		//log.Infof(context.TODO(), "found number of bundles: %n", len(bundles.Items))
		ctx.JSON(http.StatusOK, gin.H{"data": converter(bundles.Items)})
	}
}

func (bc *BundleCtrl) GetBundle(ctx *gin.Context) {

	code := ctx.Param("code")
	bundle := BundleDto{Code: code, Name: "myname"}

	ctx.JSON(http.StatusOK, gin.H{"data": bundle})
}

func converter(inBundles []v1alpha1.EntandoBundleV2) []BundleDto {
	outBundles := make([]BundleDto, 0)
	for _, bundle := range inBundles {
		b := BundleDto{
			Code:          bundle.GetAnnotations()["bundleCode"],
			Name:          bundle.GetName(),
			Title:         bundle.Spec.Title,
			Icon:          bundle.Spec.Icon,
			SignatureInfo: bundle.Spec.SignatureInfo,
		}
		outBundles = append(outBundles, b)
	}
	return outBundles
}
