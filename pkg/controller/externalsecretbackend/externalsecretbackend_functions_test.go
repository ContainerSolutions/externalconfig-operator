package externalsecretbackend

import (
	"reflect"
	"testing"

	"github.com/ContainerSolutions/externalsecret-operator/pkg/apis/externalsecretoperator/v1alpha1"
	"github.com/ContainerSolutions/externalsecret-operator/pkg/secrets"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewBackendInstanceForCR(t *testing.T) {
	Convey("Given an ExternalSecretBackend resource", t, func() {
		externalSecretBackend := v1alpha1.ExternalSecretBackend{
			Spec: v1alpha1.ExternalSecretBackendSpec{
				Type: "dummy",
				Parameters: map[string]string{
					"suffix": "-value",
				},
			},
		}
		externalSecretBackend.Name = "dummy1"
		Convey("When creating the new Backend", func() {
			err := newBackendInstanceForCR(&externalSecretBackend)
			So(err, ShouldBeNil)
			Convey("The backend is present in the backend list", func() {
				foundBackend, ok := secrets.BackendInstances["dummy1"]
				So(ok, ShouldBeTrue)
				So(reflect.TypeOf(foundBackend), ShouldEqual, reflect.TypeOf(secrets.NewDummySecretsBackend()))
				Convey("When initializing it", func() {
					err := initBackendInstanceForCR(&externalSecretBackend)
					So(err, ShouldBeNil)
					Convey("The backend is initialized correctly", func() {
						value, err := secrets.BackendInstances["dummy1"].Get("this")
						So(err, ShouldBeNil)
						So(value, ShouldEqual, "this"+"-value")
					})
				})
			})
		})
	})
}