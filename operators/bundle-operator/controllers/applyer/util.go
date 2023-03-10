package applyer

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var metadataAccessor = meta.NewAccessor()

// GetOriginalConfiguration retrieves the original configuration of the object
// from the annotation, or nil if no annotation was found.
func GetOriginalConfiguration(obj runtime.Object) ([]byte, error) {
	annots, err := metadataAccessor.Annotations(obj)
	if err != nil {
		return nil, err
	}

	if annots == nil {
		return nil, nil
	}

	original, ok := annots[corev1.LastAppliedConfigAnnotation]
	if !ok {
		return nil, nil
	}

	return []byte(original), nil
}

// GetModifiedConfiguration retrieves the modified configuration of the object.
// If annotate is true, it embeds the result as an annotation in the modified
// configuration. If an object was read from the command input, it will use that
// version of the object. Otherwise, it will use the version from the server.
func GetModifiedConfiguration(obj runtime.Object, annotate bool, codec runtime.Encoder) ([]byte, error) {
	// First serialize the object without the annotation to prevent recursion,
	// then add that serialization to it as the annotation and serialize it again.
	var modified []byte
	// Otherwise, use the server side version of the object.
	// Get the current annotations from the object.
	annots, err := metadataAccessor.Annotations(obj)
	if err != nil {
		return nil, err
	}

	if annots == nil {
		annots = map[string]string{}
	}

	original := annots[corev1.LastAppliedConfigAnnotation]
	delete(annots, corev1.LastAppliedConfigAnnotation)
	if err := metadataAccessor.SetAnnotations(obj, annots); err != nil {
		return nil, err
	}

	modified, err = runtime.Encode(codec, obj)
	if err != nil {
		return nil, err
	}

	if annotate {
		annots[corev1.LastAppliedConfigAnnotation] = string(modified)
		if err := metadataAccessor.SetAnnotations(obj, annots); err != nil {
			return nil, err
		}

		modified, err = runtime.Encode(codec, obj)
		if err != nil {
			return nil, err
		}
	}

	// Restore the object to its original condition.
	annots[corev1.LastAppliedConfigAnnotation] = original
	if err := metadataAccessor.SetAnnotations(obj, annots); err != nil {
		return nil, err
	}

	return modified, nil
}

// CreateApplyAnnotation gets the modified configuration of the object,
// without embedding it again, and then sets it on the object as the annotation.
func CreateApplyAnnotation(obj runtime.Object, codec runtime.Encoder) error {
	modified, err := GetModifiedConfiguration(obj, false, codec)
	if err != nil {
		return err
	}
	return setOriginalConfiguration(obj, modified)
}

// SetOriginalConfiguration sets the original configuration of the object
// as the annotation on the object for later use in computing a three way patch.
func setOriginalConfiguration(obj runtime.Object, original []byte) error {
	if len(original) < 1 {
		return nil
	}

	annots, err := metadataAccessor.Annotations(obj)
	if err != nil {
		return err
	}

	if annots == nil {
		annots = map[string]string{}
	}

	annots[corev1.LastAppliedConfigAnnotation] = string(original)
	return metadataAccessor.SetAnnotations(obj, annots)
}

func ConvertObjectToUnstructuredList(obj runtime.Object) ([]unstructured.Unstructured, error) {
	list := make([]unstructured.Unstructured, 0, 0)
	if meta.IsListType(obj) {
		if _, ok := obj.(*unstructured.UnstructuredList); !ok {
			return nil, fmt.Errorf("unable to convert runtime object to list")
		}

		for _, u := range obj.(*unstructured.UnstructuredList).Items {
			list = append(list, u)
		}
		return list, nil
	}

	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}

	unstructuredObj := unstructured.Unstructured{Object: unstructuredMap}
	list = append(list, unstructuredObj)

	return list, nil
}

func ConvertSingleObjectToUnstructured(obj runtime.Object) (unstructured.Unstructured, error) {
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return unstructured.Unstructured{}, err
	}

	return unstructured.Unstructured{Object: unstructuredMap}, nil
}

func MergeStringMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
