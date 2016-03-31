package mixpanel

import (
	"crypto/md5"
	"fmt"
	"strings"
	"testing"
	// "time"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	salt = "testSalt"
	path = "test/test2/?"
)

func TestMixpanelFunctions(t *testing.T) {
	//methods := []string{"test", "test2"}
	params := map[string](interface{}){
		"api_key":  "678asf8",
		"event":    []string{"pages", "signs", "test"},
		"interval": 24,
		"boolean":  false,
		"format":   "json",
		"expire":   1,
	}
	paramsJsonified := map[string]string{
		"api_key":  "678asf8",
		"event":    "[\"pages\",\"signs\",\"test\"]",
		"interval": "24",
		"boolean":  "false",
		"format":   "json",
		"expire":   "1",
	}

	// For reference, %5B and %22 correspond to ] and " respectively
	urlElements := map[string]bool{
		"api_key=678asf8": false,
		"interval=24":     false,
		"boolean=false":   false,
		"format=json":     false,
		"expire=1":        false,
		"event=%5B%22pages%22%5B%22signs%22%5B%22test%22%5B": false,
	}

	paramsHash := map[string]string{
		"c": "test1",
		"a": "test2",
		"d": "test3",
	}
	hashedURL := "a=test2c=test1d=test3" + salt
	hashedRet := fmt.Sprintf("%x", md5.Sum([]byte(hashedURL)))

	Convey("Testing mixpanel package", t, func() {
		Convey("Jsonifying function", func() {
			jsonParams, err := jsonifyParams(params)
			So(err, ShouldBeNil)
			So(jsonParams, ShouldResemble, paramsJsonified)
		})
		Convey("Testing URL forming", func() {
			encodedData, err := encodeParams(paramsJsonified)
			So(err, ShouldBeNil)

			// Checking if all the elements are present and not repeated.
			// It has to be done this way as their order is not enforced.
			uniqueElements := 0
			elements := strings.FieldsFunc(encodedData, func(r rune) bool { return r == '?' || r == '&' })
			for _, el := range elements {
				So(urlElements[el], ShouldEqual, false)
				urlElements[el] = true
				uniqueElements = uniqueElements + 1
			}
			So(uniqueElements, ShouldEqual, len(elements))
		})
		Convey("Hashing function", func() {
			hashed, _ := hashArgs(paramsHash, salt)
			So(hashed, ShouldEqual, hashedRet)
		})
	})
}
