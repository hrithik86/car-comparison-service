package controllers

import (
	"car-comparison-service/db/model"
	"car-comparison-service/db/repository/mocks"
	"car-comparison-service/errors"
	"car-comparison-service/tests"
	"car-comparison-service/utils"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestVehicleControllerTestSuite(t *testing.T) {
	suite.Run(t, new(VehicleServiceTestSuite))
}

type VehicleServiceTestSuite struct {
	suite.Suite
	mockDb        *mocks.MockIVehicle
	vehicleClient Vehicle
}

func (v *VehicleServiceTestSuite) SetupTest() {
	tests.SetupFixtures()
	ctrl := gomock.NewController(v.T())
	defer ctrl.Finish()
	v.mockDb = mocks.NewMockIVehicle(ctrl)
	v.vehicleClient = Vehicle{
		DbClient: v.mockDb,
	}
}

func (v *VehicleServiceTestSuite) TestVehicle_GetVehicleInfoById_Success() {
	Convey("Given valid uuid", v.T(), func() {
		Convey("When get by id is called", func() {
			Convey("Then it should return the data ", func() {

				vehicleData := getVehicleMockData()
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

				vehicleData := getVehicleMockData()
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

				vehicleData := getVehicleMockData()
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

func getVehicleMockData() *model.Vehicle {
	id := uuid.New()
	return &model.Vehicle{
		DbId: model.DbId{
			Id: &id,
		},
		Model:             utils.NewPtr("i30"),
		Brand:             utils.NewPtr("Hyundai"),
		ManufacturingYear: utils.NewPtr(2024),
		Type:              utils.NewPtr(model.CAR),
		Price:             utils.NewPtr(int64(4000000)),
		FuelType:          utils.NewPtr(model.PETROL),
		Mileage:           utils.NewPtr(12.4),
	}
}
