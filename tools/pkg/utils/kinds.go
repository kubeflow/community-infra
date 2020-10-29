package utils

import (
	"bytes"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

func GetObjectKind(appFile string) (string, error) {
	// Read contents
	configFileBytes, err := ioutil.ReadFile(appFile)
	if err != nil {
		return "", err
	}

	BUFSIZE := 1024
	buf := bytes.NewBufferString(string(configFileBytes))

	job := &unstructured.Unstructured{}
	err = k8syaml.NewYAMLOrJSONDecoder(buf, BUFSIZE).Decode(job)
	if err != nil {
		return "", nil
	}

	return job.GetKind(), nil
}
