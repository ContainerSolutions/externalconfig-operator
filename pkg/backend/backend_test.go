package backend

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	config "github.com/containersolutions/externalsecret-operator/pkg/config"
	. "github.com/smartystreets/goconvey/convey"
)

type MockBackend struct {
	Param1 string
}

func NewBackend() Backend {
	return &MockBackend{}
}

func (m *MockBackend) Init(params map[string]interface{}, credentials []byte) error {
	m.Param1 = params["Param1"].(string)
	return nil
}

func (m *MockBackend) Get(key string, version string) (string, error) {
	return m.Param1, nil
}

func TestRegister(t *testing.T) {
	Convey("Given a mocked backend", t, func() {
		Convey("When registering it as a backend type", func() {
			Register("mock", NewBackend)
			Convey("Then the instantiation function is registered with the correct label", func() {
				function, found := Functions["mock"]
				So(found, ShouldBeTrue)
				So(function, ShouldEqual, NewBackend)
			})
		})
	})
}

func TestInstantiate(t *testing.T) {
	Convey("Given a registered backend type", t, func() {
		Register("mock", NewBackend)
		Convey("When Instantiating it using the right label", func() {
			err := Instantiate("mock-backend", "mock")
			So(err, ShouldBeNil)
			Convey("Then a backend is instantiated with the right label", func() {
				backend, found := Instances["mock-backend"]
				So(found, ShouldBeTrue)
				So(reflect.TypeOf(backend), ShouldEqual, reflect.TypeOf(&MockBackend{}))
			})
		})
		Convey("When Instantiating it using the wrong label", func() {
			err := Instantiate("mock-backend", "mock-wrong-label")
			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unknown backend type: 'mock-wrong-label'")
			})
		})
	})
}

func TestInitFromEnv(t *testing.T) {

	configStruct := config.Config{
		Type: "mock",
		Parameters: map[string]interface{}{
			"Param1": "Value1",
		},
	}

	Convey("Given a registered backend type", t, func() {
		Register("mock", NewBackend)
		Convey("Given a valid config", func() {
			configData, _ := json.Marshal(configStruct)
			os.Setenv("OPERATOR_CONFIG", string(configData))
			Convey("When initializing backend from env", func() {
				err := InitFromEnv("mock-backend")
				So(err, ShouldBeNil)
				Convey("Then a backend is instantiated and initialized correctly", func() {
					backend, found := Instances["mock-backend"]
					So(found, ShouldBeTrue)
					So(reflect.TypeOf(backend), ShouldEqual, reflect.TypeOf(&MockBackend{}))
					value, _ := backend.Get("", "")
					So(value, ShouldEqual, "Value1")
				})
			})
		})

		Convey("Given a valid config with unknown backend type ", func() {
			configStruct.Type = "unknown"
			configData, _ := json.Marshal(configStruct)
			os.Setenv("OPERATOR_CONFIG", string(configData))
			Convey("When initializing backend from env", func() {
				err := InitFromEnv("mock-backend")
				So(err, ShouldNotBeNil)
				Convey("Then an error message is returned", func() {
					So(err.Error(), ShouldEqual, "unknown backend type: 'unknown'")
				})
			})
		})

		Convey("Given an invalid config", func() {
			os.Setenv("OPERATOR_CONFIG", "garbage")
			Convey("When initializing backend from env", func() {
				err := InitFromEnv("mock-backend")
				So(err, ShouldNotBeNil)
				Convey("Then an error is returned", func() {
					So(err.Error(), ShouldStartWith, "invalid")
				})
			})
		})

		Convey("Given a missing config", func() {
			os.Unsetenv("OPERATOR_CONFIG")
			Convey("When initializing backend from env", func() {
				err := InitFromEnv("mock-backend")
				So(err, ShouldNotBeNil)
				Convey("Then an error is returned", func() {
					So(err.Error(), ShouldStartWith, "cannot find config")
				})
			})
		})
	})
}

func TestInitFromCtrl(t *testing.T) {
	var (
		initConfig = config.Config{
			Type: "mock",
			Parameters: map[string]interface{}{
				"Param1": "Value1",
			},
			Auth: map[string]interface{}{},
		}

		credentials = `{
			Credential: "dummy-creds"
		}`
	)

	Convey("Given a registered backend type", t, func() {
		Register("mock", NewBackend)
		Convey("Given a valid config", func() {
			Convey("When initializing backend from contrl", func() {
				err := InitFromCtrl("test-ctrl", &initConfig, []byte(credentials))
				So(err, ShouldBeNil)
				Convey("Then a backend is instantiated and initialized correctly", func() {
					backend, found := Instances["test-ctrl"]
					So(found, ShouldBeTrue)
					So(reflect.TypeOf(backend), ShouldEqual, reflect.TypeOf(&MockBackend{}))
					value, _ := backend.Get("", "")
					So(value, ShouldEqual, "Value1")
				})
			})
		})

		Convey("Given a valid config with unknown backend type", func() {
			initConfig.Type = "unknown"

			Convey("When initializing backend from env", func() {
				err := InitFromCtrl("test-ctrl", &initConfig, []byte(credentials))
				So(err, ShouldNotBeNil)
				Convey("Then an error message is returned", func() {
					So(err.Error(), ShouldEqual, "unknown backend type: 'unknown'")
				})
			})
		})
	})

}
