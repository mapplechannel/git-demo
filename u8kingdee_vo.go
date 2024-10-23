package u8kingdee

import "time"

type U8TaskCenter struct {
	VrowNo        string  `gorm:"column:vrowno" json:"vrowno"`
	CbwrAstUnitId string  `gorm:"column:cbwrastunitid" json:"cbwrastunitid"`
	CbwrUnitId    string  `gorm:"column:cbwrunitid" json:"cbwrunitid"`
	NbwrAstNum    float64 `gorm:"column:nbwrastnum" json:"nbwrastnum"`
	NbwrNum       float64 `gorm:"column:nbwrnum" json:"nbwrnum"`
	CbDeptVid     string  `gorm:"column:cbdeptvid" json:"cbdeptvid"`
	CbWkId        string  `gorm:"column:cbwkid" json:"cbwkid"`
	CbTeamId      string  `gorm:"column:cbteamid" json:"cbteamid"`
	CActivityId   string  `gorm:"column:cactivityid" json:"cactivityid"`
	DTaskDate     string  `gorm:"column:dtaskdate" json:"dtaskdate"`
	NStandNum     float64 `gorm:"column:nstandnum" json:"nstandnum"`
	NActNum       float64 `gorm:"column:nactnum" json:"nactnum"`
}

type U8Department struct {
	Rank    int    `json:"rank"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	EndFlag int    `json:"endflag"`
}

type U8Supplier struct {
	BankAccNumber int64  `json:"bank_acc_number"`
	AbbrName      string `json:"abbrname"`
	Name          string `json:"name"`
	BankOpen      string `json:"bank_open"`
	Code          string `json:"code"`
	SortCode      string `json:"sort_code"`
}

type U8Unit struct {
	GroupCode string `json:"group_code"`
	Name      string `json:"name"`
	MainFlag  int    `json:"main_flag"`
	Code      string `json:"code"`
}

type U8CustomerData struct {
	AbbrName   string `json:"abbrname"`
	Name       string `json:"name"`
	DomainCode string `json:"domain_code"`
	Code       string `json:"code"`
	Industry   string `json:"industry"`
	Contact    string `json:"contact"`
	SortCode   string `json:"sort_code"`
	Mobile     int64  `json:"mobile"`
}

type U8Person struct {
	PersonGradeNo string `json:"persongradeno"`
	DepartmentNo  string `json:"departmentno"`
	Month         int    `json:"month"`
	SalaryGradeNo string `json:"salarygradeno"`
	Year          int    `json:"year"`
	PersonName    string `json:"personname"`
	PersonNo      string `json:"personno"`
}

type U8ProductionOrder struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	PkOrg        string `json:"pk_org"`
	OperatorType string `json:"operatorType"`
}

type U8Material struct {
	RecursiveFlag string `json:"recursiveflag"`
	WhCode        string `json:"whcode"`
	CinvCode      string `json:"cinvcode"`
	CinvName      string `json:"cinvname"`
	OptionalFlag  string `json:"optionalflag"`
	CwhName       string `json:"cwhname"`
	OpComponentId string `json:"opcomponentid"`
	ProductType   string `json:"producttype"`
	DrawDeptCode  string `json:"drawdeptcode"`
	CdepName      string `json:"cdepname"`
	BomId         string `json:"bomid"`
}

type U8MaterialBOM struct {
	CinvName       string    `json:"cinvname"`
	CreateUser     string    `json:"createuser"`
	Parentscrap    string    `json:"parentscrap"`
	VersionDesc    string    `json:"versiondesc"`
	BomType        string    `json:"bomtype"`
	VersionEffDate time.Time `json:"versioneffdate"`
	CinvCName      string    `json:"cinvcname"`
	BomId          string    `json:"bomid"`
	CinvCCode      string    `json:"cinvccode"`
}

// kingdee
type KingdeeDepartment struct {
	FUseOrgId    string `json:"fuseorgid"`
	FCreateOrgId string `json:"fcreateorgid"`
	FName        string `json:"fname"`
}

type KingdeeWorkCenter struct {
	FUseOrgId      string    `json:"fuseorgid"`
	FEffectDate    time.Time `json:"feffectdate"`
	FName          string    `json:"fname"`
	FCreateOrgId   string    `json:"fcreateorgid"`
	FOptCtrlCodeID string    `json:"foptctrlcodeid"`
	FDeptID        string    `json:"fdeptid"`
	FExpireDate    time.Time `json:"fexpiredate"`
}

type KingdeeSupplier struct {
	FLocAddress    string `json:"flocaddress"`
	FLocName       string `json:"flocname"`
	FLocMobile     string `json:"flocmobile"`
	FLocNewContact string `json:"flocnewcontact"`
}

type KingdeePlanOrder struct {
	FMaterialId    string    `json:"fmaterialid"`
	FInStockOrgId  string    `json:"finstockorgid"`
	FDemandOrgId   string    `json:"fdemandorgid"`
	FPlanStartDate time.Time `json:"fplanstartdate"`
	FSupplyOrgId   string    `json:"fsupplyorgid"`
}

type KingdeeUnit struct {
	FNumber         string `json:"fnumber"`
	FCurrentUnit    string `json:"fcurrentunit"`
	FOldNumber      string `json:"foldnumber"`
	FDocumentStatus string `json:"fdocumentstatus"`
	FIsBaseUnit     bool   `json:"fisbaseunit"`
	FIsSystemSet    bool   `json:"fisystemset"`
	FName           string `json:"fname"`
	FPrecision      int    `json:"fprecision"`
}

type KingdeeCustomer struct {
	FTradingCurrId string `json:"ftradingcurrid"`
	FName          string `json:"fname"`
	FCreateOrgId   string `json:"fcreateorgid"`
}

type KingdeeProductionOrder struct {
	FPrdOrgId    string    `json:"fprdorgid"`
	FOwnerTypeId string    `json:"fownertypeid"`
	FBillType    string    `json:"fbilltype"`
	FDate        time.Time `json:"fdate"`
	FPPBOMType   string    `json:"fppbomtype"`
}

type KingdeeMaterial struct {
	FCreatorId     string `json:"fcreatorid"`
	FExtVar        string `json:"fextvar"`
	FOldNumber     string `json:"foldnumber"`
	FDescription   string `json:"fdescription"`
	FMaterialGroup string `json:"fmaterialgroup"`
	FBaseProperty  string `json:"fbaseproperty"`
	FMaterialSRC   string `json:"fmaterialsrc"`
	FSpecification string `json:"fspecification"`
	FPLMMaterialId string `json:"fplmmaterialid"`
	FNumber        string `json:"fnumber"`
	FEqpClsID      string `json:"ferpclsid"`
}

type KingdeeBOM struct {
	FCreateOrgId        string `json:"fcreateorgid"`
	FMATERIALID         string `json:"fmateriaid"`
	FBOMCATEGORY        string `json:"fbomcategory"`
	FBOMUSE             string `json:"fbomuse"`
	FIsBaseUnit         bool   `json:"fisbaseunit"`
	FIsSystemSet        bool   `json:"fisystemset"`
	FName               string `json:"fname"`
	FPrecision          int    `json:"fprecision"`
	FUNITID             string `json:"funitid"`
	FMATERIALIDCHILD    string `json:"fmateriaidchild"`
	FISSUETYPE          string `json:"fissuetype"`
	FOWNERTYPEID        string `json:"fownertypeid"`
	FDOSAGETYPE         string `json:"fdosagetype"`
	FUNITIDLOT          string `json:"funitidlot"`
	FMATERIALIDLOTBASED string `json:"fmateriaidlotbased"`
	FCOBYTYPE           string `json:"fcobytype"`
	FMATERIALIDCOBY     string `json:"fmateriaidcoby"`
}

type KingdeeEmployee struct {
	FStaffNumber string    `json:"fstaffnumber"`
	FCreateOrgId string    `json:"fcreateorgid"`
	FJoinDate    time.Time `json:"fjoindate"`
	FUseOrgId    string    `json:"fuseorgid"`
	FName        string    `json:"fname"`
}
