package controllers

import (
	cacheManagerMocks "car-comparison-service/cache_manager/vehicle/mocks"
	"car-comparison-service/db/model"
	repoMocks "car-comparison-service/db/repository/mocks"
	"car-comparison-service/errors"
	"car-comparison-service/service/api/request"
	"car-comparison-service/tests"
	"car-comparison-service/utils"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestVehicleControllerTestSuite(t *testing.T) {
	suite.Run(t, new(VehicleServiceTestSuite))
}

type VehicleServiceTestSuite struct {
	suite.Suite
	mockDb                 *repoMocks.MockIVehicle
	mockVehicleSuggestions *cacheManagerMocks.MockIVehicleSuggestions
	vehicleClient          Vehicle
}

func (v *VehicleServiceTestSuite) SetupTest() {
	tests.SetupFixtures()
	ctrl := gomock.NewController(v.T())
	defer ctrl.Finish()
	v.mockDb = repoMocks.NewMockIVehicle(ctrl)
	v.mockVehicleSuggestions = cacheManagerMocks.NewMockIVehicleSuggestions(ctrl)
	v.vehicleClient = Vehicle{
		DbClient:           v.mockDb,
		VehicleSuggestions: v.mockVehicleSuggestions,
	}
}

func (v *VehicleServiceTestSuite) TestVehicle_GetVehicleInfoById_Success() {
	Convey("Given valid uuid", v.T(), func() {
		Convey("When get by id is called", func() {
			Convey("Then it should return the data ", func() {

				vehicleData := getVehicleMockData(uuid.New())
				id := vehicleData.Id
				expectedResponse := []*model.VehicleWithFeatures{{
					Vehicle:   *vehicleData,
					FeatureId: nil,
					Key:       nil,
					Value:     nil,
				}}
				v.mockDb.EXPECT().GetVehicleWithFeaturesById(gomock.Any(), *id).Times(1).Return(expectedResponse, nil)
				vehicleInfo, err := v.vehicleClient.GetVehicleInfoById(context.Background(), *id)
				So(err, ShouldBeNil)
				So(vehicleInfo, ShouldResemble, expectedResponse)
				So(vehicleInfo[0].Id, ShouldResemble, expectedResponse[0].Id)
			})
		})
	})
}

func (v *VehicleServiceTestSuite) TestVehicle_GetVehicleInfoById_Failure() {
	Convey("Given invalid uuid", v.T(), func() {
		Convey("When get by id is called", func() {
			Convey("Then it should throw error ", func() {
				id := uuid.New()
				var expectedResponse []*model.VehicleWithFeatures
				v.mockDb.EXPECT().GetVehicleWithFeaturesById(gomock.Any(), id).Times(1).Return(nil, errors.RECORD_NOT_FOUND)
				vehicleInfo, err := v.vehicleClient.GetVehicleInfoById(context.Background(), id)
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.RECORD_NOT_FOUND)
				So(vehicleInfo, ShouldResemble, expectedResponse)
			})
		})
	})
}

func (v *VehicleServiceTestSuite) TestVehicle_GetVehiclesByModelName_Success() {
	Convey("Given valid model", v.T(), func() {
		Convey("When get like model name is called", func() {
			Convey("Then it should return the data ", func() {

				vehicleData := getVehicleMockData(uuid.New())
				expectedResponse := []*model.VehicleWithAttachmentInformation{{
					Vehicle:      *vehicleData,
					AttachmentId: nil,
					MediaType:    nil,
					Path:         nil,
				}}
				v.mockDb.EXPECT().GetVehiclesByModel(gomock.Any(), utils.GetValFromPtr(vehicleData.Model)).Times(1).Return(expectedResponse, nil)
				vehicleInfo, err := v.vehicleClient.GetVehiclesByModelName(context.Background(), utils.GetValFromPtr(vehicleData.Model))
				So(err, ShouldBeNil)
				So(vehicleInfo, ShouldResemble, expectedResponse)
			})
		})
	})
}

func (v *VehicleServiceTestSuite) TestVehicle_GetVehiclesByModelName_Failure() {
	Convey("Given invalid model", v.T(), func() {
		Convey("When get like model name is called", func() {
			Convey("Then it should throw error ", func() {

				vehicleData := getVehicleMockData(uuid.New())
				var expectedResponse []*model.VehicleWithAttachmentInformation
				v.mockDb.EXPECT().GetVehiclesByModel(gomock.Any(), utils.GetValFromPtr(vehicleData.Model)).Times(1).Return(nil, errors.RECORD_NOT_FOUND)
				vehicleInfo, err := v.vehicleClient.GetVehiclesByModelName(context.Background(), utils.GetValFromPtr(vehicleData.Model))
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.RECORD_NOT_FOUND)
				So(vehicleInfo, ShouldResemble, expectedResponse)
			})
		})
	})
}

func (v *VehicleServiceTestSuite) TestVehicle_GetVehicleComparison_Success_WithHideCommonFeatures() {
	Convey("Given valid ids for comparison", v.T(), func() {
		Convey("When get vehicle comparison is called with hide common features", func() {
			Convey("Then it should return the data hiding common features", func() {

				vehicleData1 := getVehicleMockData(uuid.New())
				vehicleData2 := getVehicleMockData(uuid.New())
				ids := []uuid.UUID{*vehicleData1.Id, *vehicleData2.Id}
				expectedResponse := []*model.Vehicle{vehicleData1, vehicleData2}
				v.mockDb.EXPECT().GetVehiclesByIds(gomock.Any(), ids).Times(1).Return(expectedResponse, nil)
				respMap, err := v.vehicleClient.GetVehicleComparison(context.Background(), request.VehicleComparisonRequest{
					Ids:                ids,
					HideCommonFeatures: true,
				})
				So(err, ShouldBeNil)
				So(len(respMap), ShouldEqual, 1)
			})
		})
	})
}

func (v *VehicleServiceTestSuite) TestVehicle_GetVehicleComparison_Success_WithoutHideCommonFeatures() {
	Convey("Given valid ids for comparison", v.T(), func() {
		Convey("When get vehicle comparison is called without hiding common features", func() {
			Convey("Then it should return the data without hiding common features", func() {
				vehicleData1 := getVehicleMockData(uuid.New())
				vehicleData2 := getVehicleMockData(uuid.New())
				ids := []uuid.UUID{*vehicleData1.Id, *vehicleData2.Id}
				expectedResponse := []*model.Vehicle{vehicleData1, vehicleData2}
				v.mockDb.EXPECT().GetVehiclesByIds(gomock.Any(), ids).Times(1).Return(expectedResponse, nil)
				respMap, err := v.vehicleClient.GetVehicleComparison(context.Background(), request.VehicleComparisonRequest{
					Ids:                ids,
					HideCommonFeatures: false,
				})
				So(err, ShouldBeNil)
				assert.Greater(v.T(), len(respMap), 1)
			})
		})
	})
}

func (v *VehicleServiceTestSuite) TestVehicle_GetVehicleSuggestion_Failure() {
	Convey("Given invalid id", v.T(), func() {
		Convey("When get vehicle suggestion is called", func() {
			Convey("Then it should throw error", func() {
				vehicleData := getVehicleMockData(uuid.New())
				var expectedResponse []model.VehicleSuggestionResult
				v.mockDb.EXPECT().GetVehicleInfoById(gomock.Any(), *vehicleData.Id).Times(1).Return(nil, errors.RECORD_NOT_FOUND)
				resp, err := v.vehicleClient.GetVehicleSuggestions(context.Background(), *vehicleData.Id)
				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, errors.RECORD_NOT_FOUND)
				So(resp, ShouldResemble, expectedResponse)
			})
		})
	})
}

func (v *VehicleServiceTestSuite) TestVehicle_CreateVehicle_Success() {
	Convey("Given valid vehicle data", v.T(), func() {
		Convey("When create vehicle is called", func() {
			Convey("Then it should create vehicle", func() {
				vehicleData := getVehicleMockData(uuid.New())
				v.mockDb.EXPECT().CreateVehicle(gomock.Any(), gomock.Any()).Times(1).Return(vehicleData, nil)
				resp, err := v.vehicleClient.CreateVehicle(context.Background(), request.CreateVehicleRequest{
					Model:             vehicleData.Model,
					Brand:             vehicleData.Brand,
					ManufacturingYear: vehicleData.ManufacturingYear,
					Type:              vehicleData.Type,
					Price:             vehicleData.Price,
					FuelType:          vehicleData.FuelType,
					Mileage:           vehicleData.Mileage,
				})
				So(err, ShouldBeNil)
				So(resp, ShouldResemble, vehicleData)
			})
		})
	})
}

func (v *VehicleServiceTestSuite) TestVehicle_CreateVehicle_Failure() {
	Convey("Given invalid vehicle data", v.T(), func() {
		Convey("When create vehicle is called", func() {
			Convey("Then it should throw error", func() {
				vehicleData := getVehicleMockData(uuid.New())
				var vehicle *model.Vehicle
				v.mockDb.EXPECT().CreateVehicle(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.UNKNOWN)
				resp, err := v.vehicleClient.CreateVehicle(context.Background(), request.CreateVehicleRequest{
					Model:             nil,
					Brand:             vehicleData.Brand,
					ManufacturingYear: vehicleData.ManufacturingYear,
					Type:              vehicleData.Type,
					Price:             vehicleData.Price,
					FuelType:          vehicleData.FuelType,
					Mileage:           vehicleData.Mileage,
				})
				So(err, ShouldNotBeNil)
				So(resp, ShouldResemble, vehicle)
			})
		})
	})
}

func (v *VehicleServiceTestSuite) TestVehicle_AddVehicleAttachments_Success() {
	Convey("Given valid vehicle attachment data", v.T(), func() {
		Convey("When add vehicle attachments is called", func() {
			Convey("Then it should add vehicle attachments", func() {
				vehicleAttachments := getVehicleAttachmentsMockData()
				addVehicleAttachmentsRequest := make([]request.BulkAddVehicleAttachmentsRequest, 0, 1)
				for _, vehicleAttachment := range vehicleAttachments {
					addVehicleAttachmentsRequest = append(addVehicleAttachmentsRequest, request.BulkAddVehicleAttachmentsRequest{
						Name:      vehicleAttachment.Name,
						Path:      vehicleAttachment.Path,
						MediaType: vehicleAttachment.MediaType,
					})
				}
				v.mockDb.EXPECT().BulkAddVehicleAttachments(gomock.Any(), gomock.Any()).Times(1).Return(vehicleAttachments, nil)
				resp, err := v.vehicleClient.AddVehicleAttachments(context.Background(), *vehicleAttachments[0].VehicleId, addVehicleAttachmentsRequest)
				So(err, ShouldBeNil)
				So(resp, ShouldResemble, vehicleAttachments)
			})
		})
	})
}

func (v *VehicleServiceTestSuite) TestVehicle_AddVehicleAttachments_Failure() {
	Convey("Given invalid vehicle attachment data", v.T(), func() {
		Convey("When add vehicle attachments is called", func() {
			Convey("Then it should throw error", func() {
				vehicleAttachments := getVehicleAttachmentsMockData()
				var expectedResponse []*model.VehicleAttachment
				addVehicleAttachmentsRequest := make([]request.BulkAddVehicleAttachmentsRequest, 0, 1)
				for _, vehicleAttachment := range vehicleAttachments {
					addVehicleAttachmentsRequest = append(addVehicleAttachmentsRequest, request.BulkAddVehicleAttachmentsRequest{
						Name:      nil,
						Path:      vehicleAttachment.Path,
						MediaType: vehicleAttachment.MediaType,
					})
				}
				v.mockDb.EXPECT().BulkAddVehicleAttachments(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.UNKNOWN)
				resp, err := v.vehicleClient.AddVehicleAttachments(context.Background(), *vehicleAttachments[0].VehicleId, addVehicleAttachmentsRequest)
				So(err, ShouldNotBeNil)
				So(resp, ShouldResemble, expectedResponse)
			})
		})
	})
}

func (v *VehicleServiceTestSuite) TestVehicle_AddVehicleFeatures_Success() {
	Convey("Given valid vehicle features data", v.T(), func() {
		Convey("When add vehicle features is called", func() {
			Convey("Then it should add vehicle features", func() {
				vehicleFeatures := getVehicleFeaturesMockData()
				addVehicleFeaturesRequest := make([]request.BulkAddVehicleFeaturesRequest, 0, 1)
				for _, vehicleFeature := range vehicleFeatures {
					addVehicleFeaturesRequest = append(addVehicleFeaturesRequest, request.BulkAddVehicleFeaturesRequest{
						Key:   vehicleFeature.Key,
						Value: vehicleFeature.Value,
					})
				}
				v.mockDb.EXPECT().BulkAddVehicleFeatures(gomock.Any(), gomock.Any()).Times(1).Return(vehicleFeatures, nil)
				resp, err := v.vehicleClient.AddVehicleFeatures(context.Background(), *vehicleFeatures[0].VehicleId, addVehicleFeaturesRequest)
				So(err, ShouldBeNil)
				So(resp, ShouldResemble, vehicleFeatures)
			})
		})
	})
}

func (v *VehicleServiceTestSuite) TestVehicle_AddVehicleFeatures_Failure() {
	Convey("Given invalid vehicle features data", v.T(), func() {
		Convey("When add vehicle features is called", func() {
			Convey("Then it should throw error", func() {
				vehicleFeatures := getVehicleFeaturesMockData()
				var expectedResponse []*model.VehicleFeatures
				addVehicleFeaturesRequest := make([]request.BulkAddVehicleFeaturesRequest, 0, 1)
				for _, vehicleFeature := range vehicleFeatures {
					addVehicleFeaturesRequest = append(addVehicleFeaturesRequest, request.BulkAddVehicleFeaturesRequest{
						Key:   nil,
						Value: vehicleFeature.Value,
					})
				}
				v.mockDb.EXPECT().BulkAddVehicleFeatures(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.UNKNOWN)
				resp, err := v.vehicleClient.AddVehicleFeatures(context.Background(), *vehicleFeatures[0].VehicleId, addVehicleFeaturesRequest)
				So(err, ShouldNotBeNil)
				So(resp, ShouldResemble, expectedResponse)
			})
		})
	})
}

func (v *VehicleServiceTestSuite) TestVehicle_GetVehicleSuggestion_CachedResultSuccess() {
	Convey("Given valid vehicle id", v.T(), func() {
		Convey("When get vehicle suggestion is called", func() {
			Convey("Then it should returned cached result", func() {
				vehicleId := uuid.New()
				key := fmt.Sprintf("{suggestions-set}_%s", vehicleId.String())
				vehicleInfo := getVehicleMockData(vehicleId)
				suggestionsData := getVehicleSuggestionsMockData()
				v.mockVehicleSuggestions.EXPECT().CreateKey(gomock.Any(), gomock.Any()).Times(1).Return(key)
				v.mockVehicleSuggestions.EXPECT().GetVehicleSuggestionsDetails(gomock.Any(), key).Times(1).Return(suggestionsData, nil)
				v.mockDb.EXPECT().GetVehicleInfoById(gomock.Any(), gomock.Any()).Times(1).Return(vehicleInfo, nil)
				resp, err := v.vehicleClient.GetVehicleSuggestions(context.Background(), uuid.New())
				So(err, ShouldBeNil)
				So(resp, ShouldResemble, suggestionsData)
			})
		})
	})
}

func getVehicleMockData(id uuid.UUID) *model.Vehicle {
	return &model.Vehicle{
		DbId: model.DbId{
			Id: &id,
		},
		Model:             utils.NewPtr("i20"),
		Brand:             utils.NewPtr("Hyundai"),
		ManufacturingYear: utils.NewPtr(2024),
		Type:              utils.NewPtr(model.CAR),
		Price:             utils.NewPtr(int64(4000000)),
		FuelType:          utils.NewPtr(model.PETROL),
		Mileage:           utils.NewPtr(12.4),
	}
}

func getVehicleAttachmentsMockData() []*model.VehicleAttachment {
	return []*model.VehicleAttachment{{
		DbId: model.DbId{
			Id: utils.NewPtr(uuid.New()),
		},
		Name:      utils.NewPtr("Img-1"),
		Path:      utils.NewPtr("s3://images/Img-1.jpg"),
		MediaType: utils.NewPtr(model.IMAGE),
		VehicleId: utils.NewPtr(uuid.New()),
	},
	}
}

func getVehicleFeaturesMockData() []*model.VehicleFeatures {
	return []*model.VehicleFeatures{{
		DbId: model.DbId{
			Id: utils.NewPtr(uuid.New()),
		},
		Key:       utils.NewPtr("Power Steering"),
		Value:     utils.NewPtr("true"),
		VehicleId: utils.NewPtr(uuid.New()),
	},
	}
}

func getVehicleSuggestionsMockData() []model.VehicleSuggestionResult {
	return []model.VehicleSuggestionResult{
		{
			Id:                uuid.New(),
			Model:             "i10",
			Brand:             "Hyundai",
			ManufacturingYear: 2022,
			Price:             1000000,
			Mileage:           20.2,
			FuelType:          string(model.DIESEL),
			Type:              string(model.CAR),
			Rank:              1,
		},
		{
			Id:                uuid.New(),
			Model:             "Grand-i10",
			Brand:             "Hyundai",
			ManufacturingYear: 2023,
			Price:             1500000,
			Mileage:           20.2,
			FuelType:          string(model.DIESEL),
			Type:              string(model.CAR),
			Rank:              2,
		},
	}
}
