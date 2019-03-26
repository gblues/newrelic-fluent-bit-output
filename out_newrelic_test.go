package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Out New Relic", func() {

	Describe("Prepares payload", func() {
		gomega.RegisterFailHandler(ginkgo.Fail)
		It("converts the map[interface{}] inteface{} to map[string] interface[], "+
			"updates the timestamp, and renames the log field to message",
			func() {
				inputMap := make(map[interface{}]interface{})
				var inputTimestamp interface{}
				inputTimestamp = output.FLBTime{
					time.Now(),
				}
				inputMap["log"] = "message"
				foundOutput := prepareRecord(inputMap, inputTimestamp)
				fmt.Println(foundOutput)
				Expect(foundOutput["message"]).To(Equal("message"))
				Expect(foundOutput["log"]).To(BeNil())
				Expect(foundOutput["timestamp"]).To(Equal(inputTimestamp.(output.FLBTime).UnixNano() / 1000000))
			})
	})

	Describe("HTTP Request body", func() {
		gomega.RegisterFailHandler(ginkgo.Fail)
		var server *ghttp.Server
		const expectedInsertKey = "sweetKey"
		var expectedEndpoint string
		BeforeEach(func() {
			server = ghttp.NewServer()
			expectedEndpoint = server.URL() + "/v1/logs"
			server.AppendHandlers(
				ghttp.CombineHandlers(ghttp.VerifyRequest("POST", "/v1/logs"),
					ghttp.VerifyHeader(http.Header{
						"X-Insert-Key":     []string{expectedInsertKey},
						"Content-Type":     []string{"application/json"},
						"Content-Encoding": []string{"gzip"},
					}),
				))
		})

		AfterEach(func() {
			server.Close()
		})
		It("correctly packages and posts json", func() {

			var testConfig = PluginConfig{
				apiKey:        expectedInsertKey,
				endpoint:      expectedEndpoint,
				maxBufferSize: 256000,
				maxRecords:    2,
			}

			var testRecords []map[string]interface{}
			var testRecord map[string]interface{}

			testRecord = make(map[string]interface{})
			testRecord["timestamp"] = time.Now().UnixNano() / int64(time.Millisecond)
			testRecord["message"] = "cool story"
			testRecords = append(testRecords, testRecord)
			responseChan := prepare(testRecords, &testConfig)
			<-responseChan
		})
	})
})