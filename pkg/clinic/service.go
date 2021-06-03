package clinic

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"go.uber.org/zap"
)

type DataDownloader struct {
	logger *zap.Logger
}

func NewDataDownloader(l *zap.Logger) *DataDownloader {
	return &DataDownloader{
		logger: l,
	}
}

func (d *DataDownloader) GetClinicData(logger *zap.Logger) ([]Clinic, error) {
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
