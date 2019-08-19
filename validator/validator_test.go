package validator

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update golder files")

func ExampleValidator_Run() {
	const cache = "../testdata/web3_cache"
	const file = "000007700162.json"
	var dirs = []string{
		"eth_getBlockByNumber",
		"eth_getTransactionReceipt",
		"eth_getUncleByBlockHashAndIndex",
		"trace_block",
		"trace_replayBlockTransactions",
	}

	v := New()

	for _, dir := range dirs {
		path := cache + "/" + dir + "/" + file

		// Make sure the file exists before trying to read it
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			log.Fatal("could not read file: ", err)
		}

		data, err := ioutil.ReadFile(cache + "/" + dir + "/" + file)
		if err != nil {
			log.Fatal("could not read file: ", err)
		}

		switch dir {
		case "eth_getBlockByNumber":
			err = v.LoadBlockResponse(data)
		case "eth_getTransactionReceipt":
			err = v.LoadReceiptsResponse(data)
		case "eth_getUncleByBlockHashAndIndex":
			err = v.LoadUnclesResponse(data)
		case "trace_block":
			err = v.LoadTraceBlockResponse(data)
		case "trace_replayBlockTransactions":
			err = v.LoadReplayResponse(data)
		}

		if err != nil {
			log.Fatal("could not load data into validator: ", err)
		}
	}

	ok, err := v.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Validator response for block %s: %s", file, func(valid bool) string {
		if valid {
			return "valid"
		} else {
			return "not valid"
		}
	}(ok))

	// Output: Validator response for block 000007700162.json: valid
}

func TestChecker_Load(t *testing.T) {
	const cache = "../testdata/web3_cache"
	const golden = "../testdata/golden"
	var dirs = []string{
		"eth_getBlockByNumber",
		"eth_getTransactionReceipt",
		"eth_getUncleByBlockHashAndIndex",
		"trace_block",
		"trace_replayBlockTransactions",
	}

	files, err := ioutil.ReadDir(cache + "/eth_getBlockByNumber")
	if err != nil {
		t.Error(err)
	}

	for _, f := range files {
		t.Run(f.Name(), func(tt *testing.T) {
			v := New()

			for _, dir := range dirs {
				path := cache + "/" + dir + "/" + f.Name()

				// Make sure the file exists before trying to read it
				_, err := os.Stat(path)
				if os.IsNotExist(err) {
					continue
				} else if err != nil {
					tt.Error(err)
				}

				data, err := ioutil.ReadFile(cache + "/" + dir + "/" + f.Name())
				if err != nil {
					tt.Error(err)
				}

				switch dir {
				case "eth_getBlockByNumber":
					err = v.LoadBlockResponse(data)
				case "eth_getTransactionReceipt":
					err = v.LoadReceiptsResponse(data)
				case "eth_getUncleByBlockHashAndIndex":
					err = v.LoadUnclesResponse(data)
				case "trace_block":
					err = v.LoadTraceBlockResponse(data)
				case "trace_replayBlockTransactions":
					err = v.LoadReplayResponse(data)
				}

				if err != nil {
					tt.Error(err)
				}
			}

			checkerJSON, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				tt.Error(err)
			}

			if *update {
				err = ioutil.WriteFile(golden+"/"+f.Name(), checkerJSON, 0644)
				if err != nil {
					t.Fatal("could not write golden file")
				}

				return
			}

			expectedOutput, err := ioutil.ReadFile(golden + "/" + f.Name())
			if err != nil {
				t.Fatal("got error reading golden file ", f.Name(), ":", err)
			}

			var expectedChecker *Validator
			err = json.Unmarshal(expectedOutput, &expectedChecker)
			if err != nil {
				t.Error("got error decoding expected output ", f.Name(), ":", err)
			}

			assert.Equal(tt, expectedChecker, v)
		})
	}
}

func TestChecker_Verify(t *testing.T) {
	const cache = "../testdata/web3_cache"
	var dirs = []string{
		"eth_getBlockByNumber",
		"eth_getTransactionReceipt",
		"eth_getUncleByBlockHashAndIndex",
		"trace_block",
		"trace_replayBlockTransactions",
	}

	files, err := ioutil.ReadDir(cache + "/eth_getBlockByNumber")
	if err != nil {
		t.Error(err)
	}

	for _, f := range files {
		t.Run(f.Name(), func(tt *testing.T) {
			v := New()

			for _, dir := range dirs {
				path := cache + "/" + dir + "/" + f.Name()

				// Make sure the file exists before trying to read it
				_, err := os.Stat(path)
				if os.IsNotExist(err) {
					continue
				} else if err != nil {
					tt.Error(err)
				}

				data, err := ioutil.ReadFile(cache + "/" + dir + "/" + f.Name())
				if err != nil {
					tt.Error(err)
				}

				switch dir {
				case "eth_getBlockByNumber":
					err = v.LoadBlockResponse(data)
				case "eth_getTransactionReceipt":
					err = v.LoadReceiptsResponse(data)
				case "eth_getUncleByBlockHashAndIndex":
					err = v.LoadUnclesResponse(data)
				case "trace_block":
					err = v.LoadTraceBlockResponse(data)
				case "trace_replayBlockTransactions":
					err = v.LoadReplayResponse(data)
				}

				if err != nil {
					tt.Error(err)
				}
			}

			ok, err := v.Run()
			if err != nil {
				tt.Error(err)
			}

			assert.True(tt, ok)
		})
	}
}

func BenchmarkValidator_RunSmall(b *testing.B) {
	benchmarkRun("000007700162.json", b)
}

func BenchmarkValidator_RunMedium(b *testing.B) {
	benchmarkRun("000007714301.json", b)
}

func BenchmarkValidator_RunBig(b *testing.B) {
	benchmarkRun("000007000062.json", b)
}

func benchmarkRun(file string, b *testing.B) {
	const cache = "../testdata/web3_cache"
	var dirs = []string{
		"eth_getBlockByNumber",
		"eth_getTransactionReceipt",
		"eth_getUncleByBlockHashAndIndex",
		"trace_block",
		"trace_replayBlockTransactions",
	}

	v := New()

	for _, dir := range dirs {
		path := cache + "/" + dir + "/" + file

		// Make sure the file exists before trying to read it
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			b.Error(err)
		}

		data, err := ioutil.ReadFile(cache + "/" + dir + "/" + file)
		if err != nil {
			b.Error(err)
		}

		switch dir {
		case "eth_getBlockByNumber":
			err = v.LoadBlockResponse(data)
		case "eth_getTransactionReceipt":
			err = v.LoadReceiptsResponse(data)
		case "eth_getUncleByBlockHashAndIndex":
			err = v.LoadUnclesResponse(data)
		case "trace_block":
			err = v.LoadTraceBlockResponse(data)
		case "trace_replayBlockTransactions":
			err = v.LoadReplayResponse(data)
		}

		if err != nil {
			b.Error(err)
		}
	}

	for i := 0; i < b.N; i++ {
		v.Run()
	}
}
