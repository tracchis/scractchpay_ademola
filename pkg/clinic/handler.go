package clinic

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/scratchpay_ademola/internal/httputil"
	"github.com/scratchpay_ademola/internal/logger"
	"github.com/scratchpay_ademola/internal/validatorutil"

	"github.com/thedevsaddam/gojsonq/v2"
	"go.uber.org/zap"
)

const (
	vetClinicsURL    = "https://storage.googleapis.com/scratchpay-code-challenge/vet-clinics.json"
	dentalClinicsURL = "https://storage.googleapis.com/scratchpay-code-challenge/dental-clinics.json"
)

func Search() http.HandlerFunc {
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
			l.Error("Failed parsing json to user struct", zap.Error(err))
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

		data, err := getClinicData(l)
		if err != nil {
			l.Error("error fetching clinic data")
			httputil.JSONError(w, http.StatusBadRequest, "error fetching all clinics", attrErrMessages)
		}

		d, err := json.Marshal(data)
		query := gojsonq.New().
			FromString(string(d)).
			WhereContains("name", params.Name).
			Where("state", skipEmpty(params.State), params.State).
			Where("availability.from", skipEmpty(params.From), params.From).
			Where("availability.to", skipEmpty(params.To), params.To)

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

func skipEmpty(query string) string {
	if query == "" {
		return "!="
	}

	return "="
}

func GetAllClinics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.From(context.Background())
		attrErrMessages := validatorutil.GetAttributeErrorMessages()

		data, err := getClinicData(l)
		if err != nil {
			l.Error("error fetching clinic data")
			httputil.JSONError(w, http.StatusBadRequest, "error fetching all clinics", attrErrMessages)
		}

		httputil.JSONSuccess(w, http.StatusOK, data)
	}
}

func fetchData(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getDentalClinics() ([]Clinic, error) {
	body, err := fetchData(dentalClinicsURL)
	if err != nil {
		return nil, err
	}

	var clinics []DentalClinic
	err = json.Unmarshal(body, &clinics)
	if err != nil {
		return nil, err
	}

	return parseDentalClinics(clinics), nil
}

func parseDentalClinics(dentalClinics []DentalClinic) []Clinic {
	var clinics []Clinic
	for _, cl := range dentalClinics {
		clinics = append(clinics, Clinic{
			Name:         cl.Name,
			State:        cl.State,
			Availability: cl.Availability,
		})
	}

	return clinics
}

func getVetClinics() ([]Clinic, error) {

	body, err := fetchData(vetClinicsURL)
	if err != nil {
		return nil, err
	}

	var clinics []VetClinic
	err = json.Unmarshal(body, &clinics)
	if err != nil {
		return nil, err
	}

	return parseVetClinics(clinics), nil
}

func parseVetClinics(vetClinics []VetClinic) []Clinic {
	var clinics []Clinic
	for _, cl := range vetClinics {
		clinics = append(clinics, Clinic{
			Name:         cl.Name,
			State:        cl.State,
			Availability: cl.Availability,
		})
	}

	return clinics
}

func getClinicData(logger *zap.Logger) ([]Clinic, error) {
	var clinics []Clinic

	var sg sync.WaitGroup
	sg.Add(2)

	go func() {
		dentalClinics, err := getDentalClinics()
		if err != nil {
			logger.Error("error fetching dental clinics", zap.Error(err))
		}
		clinics = append(clinics, dentalClinics...)
		sg.Done()
	}()

	go func() {
		vetClinics, err := getVetClinics()
		if err != nil {
			logger.Error("error fetching vet clinics", zap.Error(err))
		}
		clinics = append(clinics, vetClinics...)
		sg.Done()
	}()

	sg.Wait()
	return clinics, nil
}
