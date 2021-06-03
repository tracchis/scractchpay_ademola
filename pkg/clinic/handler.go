package clinic

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/scratchpay_ademola/internal/httputil"
	"github.com/scratchpay_ademola/internal/logger"
	"github.com/scratchpay_ademola/internal/validatorutil"

	"fmt"
	"time"

	"github.com/thedevsaddam/gojsonq/v2"
	"go.uber.org/zap"
)

const (
	vetClinicsURL    = "https://storage.googleapis.com/scratchpay-code-challenge/vet-clinics.json"
	dentalClinicsURL = "https://storage.googleapis.com/scratchpay-code-challenge/dental-clinics.json"
)

type DataFetcher interface {
	GetClinicData(logger *zap.Logger) ([]Clinic, error)
}

func GetAllClinics(dataFetcher DataFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.From(context.Background())
		attrErrMessages := validatorutil.GetAttributeErrorMessages()

		data, err := dataFetcher.GetClinicData(l)
		if err != nil {
			l.Error("error fetching clinic data")
			httputil.JSONError(w, http.StatusInternalServerError, "error fetching all clinics", attrErrMessages)
			return
		}

		httputil.JSONSuccess(w, http.StatusOK, data)
	}
}

func Search(fetcher DataFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.From(context.Background())
		attrErrMessages := validatorutil.GetAttributeErrorMessages()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			l.Error("failed fetching json body")
			httputil.JSONError(w, http.StatusBadRequest, "internal server error", attrErrMessages)
			return
		}

		var params SearchParams
		err = json.Unmarshal(body, &params)
		if err != nil {
			l.Error("Failed parsing json to params struct", zap.Error(err))
			httputil.JSONError(w, http.StatusBadRequest, "invalid json params", attrErrMessages)
			return
		}

		validate := validatorutil.GetValidator()

		err = validate.Struct(params)
		if err != nil {
			attrErrMessages = validatorutil.GetTranslatedErrors(err)
			httputil.JSONError(w, http.StatusBadRequest, "invalid attributes", attrErrMessages)
			return
		}

		data, err := fetcher.GetClinicData(l)
		if err != nil {
			l.Error("error fetching clinic data")
			httputil.JSONError(w, http.StatusInternalServerError, "error fetching all clinics", attrErrMessages)
			return
		}

		d, err := json.Marshal(data)
		query := gojsonq.New().
			FromString(string(d))

		query.Macro("date<=", dateLessOrEqualTo)
		query.Macro("date>=", dateGreaterOrEqualTo)

		if params.Name != "" {
			query.WhereContains("name", params.Name)
		}

		if params.State != "" {
			query.WhereContains("state", params.State)
		}

		if params.To != "" {
			query.Where("availability.from", "date>=", params.From)
		}

		if params.From != "" {
			query.Where("availability.to", "date<=", params.To)
		}

		result := query.Get()

		var clinics []Clinic
		err = mapstructure.Decode(result, &clinics)
		if err != nil {
			l.Error("failed decoding clinics", zap.Error(err))
			httputil.JSONError(w, http.StatusInternalServerError, "error searching clinics", attrErrMessages)
			return
		}

		httputil.JSONSuccess(w, http.StatusOK, clinics)
	}
}

const layout = "2006-01-02"

func dateLessOrEqualTo(x, y interface{}) (bool, error) {
	xs, okx := x.(string)
	ys, oky := y.(string)
	if !okx || !oky {
		return false, fmt.Errorf("date support for string only")
	}

	t1, _ := time.Parse(layout, xs)
	t2, _ := time.Parse(layout, ys)

	return t1.Unix() <= t2.Unix(), nil
}

func dateGreaterOrEqualTo(x, y interface{}) (bool, error) {
	xs, okx := x.(string)
	ys, oky := y.(string)
	if !okx || !oky {
		return false, fmt.Errorf("date support for string only")
	}

	t1, _ := time.Parse(layout, xs)
	t2, _ := time.Parse(layout, ys)

	return t1.Unix() >= t2.Unix(), nil
}
