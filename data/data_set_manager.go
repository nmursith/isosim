package data

import (
	"errors"
	"github.com/rkbalgi/isosim/web/spec"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"encoding/json"
	"bytes"
	"github.com/rkbalgi/isosim/web/ui_data"
	"strconv"
	"strings"
)

type dataSetManager struct{}

var instance *dataSetManager
var dataDir string

func Init(dirname string) error {
	dir, err := os.Open(dirname)
	if err != nil {
		return err
	}
	dir.Close()
	dataDir = dirname
	return nil
}

func DataSetManager() *dataSetManager {

	init := sync.Once{}
	init.Do(func() {
		instance = new(dataSetManager)

	})
	return instance
}

var DataSetExistsError = errors.New("Data Set Exists.")

func checkIfExists(specId string, msgId string, name string) (bool, error) {

	//check if the dir exists for this spec and msg
	//and if not create one first

	dir, err := os.Open(filepath.Join(dataDir, specId, msgId))
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join(dataDir, specId, msgId), os.ModeDir)
		if err != nil {
			return false, err
		}
		dir, err = os.Open(filepath.Join(dataDir, specId, msgId))
		if err != nil {
			return false, err
		}

	}

	fiSlice, err := dir.Readdir(-1)
	if err != nil {
		return false, err
	}
	for _, fi := range fiSlice {
		if fi.Name() == name {
			return true, nil
		}
	}

	return false, nil

}

//Returns a list of all data sets (names only) for the given specId
//and msgId
func (dsm *dataSetManager) GetAll(specId string, msgId string) ([]string, error) {

	dir, err := os.Open(filepath.Join(dataDir, specId, msgId))
	if err != nil {
		return nil, err

	}

	fi, err := dir.Readdir(-1)

	if err != nil {
		return nil, err
	}

	var dataSets = make([]string, 0, 10)
	for _, ds := range fi {

		if !ds.IsDir() {
			dataSets = append(dataSets, ds.Name())
		}
	}

	return dataSets, nil

}

//Returns the content of a specific data set
func (dsm *dataSetManager) Get(specId string, msgId string, dsName string) ([]byte, error) {

	//file,err:=os.Open(filepath.Join(dataDir,specId,msgId,dsName));
	//if err!=nil{
	//	//	return nil,err;

	//	}
	data, err := ioutil.ReadFile(filepath.Join(dataDir, specId, msgId, dsName))
	if err != nil {
		return nil, err

	}
	return data, nil

}

func (dsm *dataSetManager) Add(specId string, msgId string, name string, data string) error {

	if spec.DebugEnabled {
		log.Print("Adding data set - " + name + " data = " + data)
	}

	exists, err := checkIfExists(specId, msgId, name)
	if err != nil {
		return err
	}
	if exists {
		return DataSetExistsError
	}

	err = ioutil.WriteFile(filepath.Join(dataDir, specId, msgId, name), []byte(data), os.FileMode(os.O_CREATE))
	if err != nil {
		return err
	}
	return nil
}

func (dsm *dataSetManager) AddServerDef(defString string) (string,error) {

	if spec.DebugEnabled {
		log.Print("Adding server definition - .. JSON = "+defString);
	}

	serverDef:=&ui_data.ServerDef{};
	err:=json.NewDecoder(bytes.NewBufferString(defString)).Decode(serverDef);
	if(err!=nil){
		return "",err;
	}


	strSpecId:=strconv.Itoa(serverDef.SpecId);
	dir,err:=os.Open(filepath.Join(dataDir,strSpecId));
	if err!=nil && os.IsNotExist(err){
		//create dir if one doesn't exist
		os.Mkdir(filepath.Join(dataDir,strSpecId),os.ModeDir)
	}else{
		if err!=nil{
			return "", err;
		}
	}
	dir.Close();

	fileName:=serverDef.ServerName;
	fileName=strings.Replace(fileName," ","",-1);
	fileName=strings.Replace(fileName,",","",-1);
	fileName=fileName+".srvdef.json";

	log.Print("Writing spec def to file = "+fileName);

	defFile,err:=os.Open(filepath.Join(dataDir,strSpecId,fileName));
	if err==nil{
		return "",errors.New("Def file exists. Please use a different name.");
	}
	defFile.Close();

	err = ioutil.WriteFile(filepath.Join(dataDir, strSpecId,fileName), []byte(defString), os.FileMode(os.O_CREATE))
	if err != nil {
		return "",err
	}
	return fileName,nil
}

func (dm *dataSetManager) GetServerDefs(specId string) ([]string,error){

	dir,err:=os.Open(filepath.Join(dataDir,specId));
	if err!=nil{
		return nil,err;
	}
	fi,err:=dir.Readdir(-1);
	res:=make([]string,0,2);
	for _,fi:=range(fi){
		if strings.HasSuffix(fi.Name(),".srvdef.json"){
			res=append(res,fi.Name());
		}
	}

	return res,nil;


}


func (dm *dataSetManager) GetServerDef(specId string,name string) ([]byte,error){

	file,err:=os.Open(filepath.Join(dataDir,specId,name));
	if err!=nil{
		return nil,err;
	}
	data,err:=ioutil.ReadAll(file);
	return data,err;


}

func (dsm *dataSetManager) Update(specId string, msgId string, name string, data string) error {

	if spec.DebugEnabled {
		log.Print("Updating data set - " + name + " data = " + data)
	}

	err := ioutil.WriteFile(filepath.Join(dataDir, specId, msgId, name), []byte(data), os.FileMode(os.O_TRUNC))
	if err != nil {
		return err
	}
	return nil
}


