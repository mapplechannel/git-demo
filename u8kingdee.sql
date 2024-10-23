CREATE TABLE u8_task_center (
    vrowno VARCHAR(10),              
    cbwrastunitid VARCHAR(10),       
    cbwrunitid VARCHAR(10),           
    nbwrastnum NUMERIC(10, 2),       
    nbwrnum NUMERIC(10, 2),          
    cbdeptvid VARCHAR(10),            
    cbwkid VARCHAR(50),              
    cbteamid VARCHAR(10),             
    cactivityid VARCHAR(10),          
    dtaskdate TIMESTAMP,             
    nstandnum NUMERIC(10, 2),         
    nactnum NUMERIC(10, 2)            
);

CREATE TABLE u8_department (
    rank INT,
    name VARCHAR(255),
    code VARCHAR(10),
    endflag INT
);

CREATE TABLE u8_supplier (
    bank_acc_number BIGINT,
    abbrname VARCHAR(255),
    name VARCHAR(255),
    bank_open VARCHAR(255),
    code VARCHAR(10),
    sort_code VARCHAR(10)
);

CREATE TABLE u8_unit (
    group_code VARCHAR(10),
    name VARCHAR(255),
    main_flag INT,
    code VARCHAR(10)
);

CREATE TABLE u8_customer_data (
    abbrname VARCHAR(255),
    name VARCHAR(255),
    domain_code VARCHAR(10),
    code VARCHAR(10),
    industry VARCHAR(50),
    contact VARCHAR(50),
    sort_code VARCHAR(10),
    mobile BIGINT
);

CREATE TABLE u8_person (
    persongradeno VARCHAR(50),
    departmentno VARCHAR(10),
    month INT,
    salarygradeno VARCHAR(10),
    year INT,
    personname VARCHAR(255),
    personno VARCHAR(50)
);

CREATE TABLE u8_production_order (
    code VARCHAR(50),
    name VARCHAR(255),
    pk_org VARCHAR(50),
    operator_type VARCHAR(50)
);

CREATE TABLE u8_material (
    recursiveflag VARCHAR(10),
    whcode VARCHAR(10),
    cinvcode VARCHAR(50),
    cinvname VARCHAR(255),
    optionalflag VARCHAR(10),
    cwhname VARCHAR(255),
    opcomponentid VARCHAR(50),
    producttype VARCHAR(10),
    drawdeptcode VARCHAR(10),
    cdepname VARCHAR(255),
    bomid VARCHAR(50)
);

CREATE TABLE u8_material_bom (
    cinvname VARCHAR(255),
    createuser VARCHAR(50),
    parentscrap VARCHAR(10),
    versiondesc VARCHAR(10),
    bomtype VARCHAR(10),
    versioneffdate TIMESTAMP,
    cinvcname VARCHAR(255),
    bomid VARCHAR(50),
    cinvccode VARCHAR(10)
);

CREATE TABLE kingdee_department (
    fuseorgid VARCHAR(50),
    fcreateorgid VARCHAR(50),
    fname VARCHAR(100)
);

CREATE TABLE kingdee_work_center (
    fuseorgid VARCHAR(50),
    feffectdate TIMESTAMP,
    fname VARCHAR(100),
    fcreateorgid VARCHAR(50),
    foptctrlcodeid VARCHAR(50),
    fdeptid VARCHAR(50),
    fexpiredate TIMESTAMP
);

CREATE TABLE kingdee_supplier (
    flocaddress VARCHAR(255),
    flocname VARCHAR(100),
    flocmobile VARCHAR(20),
    flocnewcontact VARCHAR(100)
);

CREATE TABLE kingdee_plan_order (
    fmaterialid VARCHAR(50),
    finstockorgid VARCHAR(50),
    fdemandorgid VARCHAR(50),
    fplanstartdate TIMESTAMP,
    fsupplyorgid VARCHAR(50)
);

CREATE TABLE kingdee_unit (
    fnumber VARCHAR(50),
    fcurrentunit VARCHAR(50),
    foldnumber VARCHAR(50),
    fdocumentstatus VARCHAR(100),
    fisbaseunit BOOLEAN,
    fisystemset BOOLEAN,
    fname VARCHAR(100),
    fprecision INT
);

CREATE TABLE kingdee_customer (
    ftradingcurrid VARCHAR(50),
    fname VARCHAR(100),
    fcreateorgid VARCHAR(50)
);

CREATE TABLE kingdee_production_order (
    fprdorgid VARCHAR(50),
    fownertypeid VARCHAR(50),
    fbilltype VARCHAR(50),
    fdate TIMESTAMP,
    fppbomtype VARCHAR(50)
);

CREATE TABLE kingdee_material (
    fcreatorid VARCHAR(50),
    fextvar VARCHAR(50),
    foldnumber VARCHAR(50),
    fdescription VARCHAR(255),
    fmaterialgroup VARCHAR(50),
    fbaseproperty VARCHAR(50),
    fmaterialsrc VARCHAR(50),
    fspecification VARCHAR(100),
    fplmmaterialid VARCHAR(50),
    fnumber VARCHAR(50),
    ferpclsid VARCHAR(50)
);

CREATE TABLE kingdee_bom (
    fcreateorgid VARCHAR(50),
    fmateriaid VARCHAR(50),
    fbomcategory VARCHAR(50),
    fbomuse VARCHAR(50),
    fisbaseunit BOOLEAN,
    fisystemset BOOLEAN,
    fname VARCHAR(100),
    fprecision INT,
    funitid VARCHAR(50),
    fmateriaidchild VARCHAR(50),
    fissuetype VARCHAR(50),
    fownertypeid VARCHAR(50),
    fdosagetype VARCHAR(50),
    funitidlot VARCHAR(50),
    fmateriaidlotbased VARCHAR(50),
    fcobytype VARCHAR(50),
    fmateriaidcoby VARCHAR(50)
);

CREATE TABLE kingdee_employee (
    fstaffnumber VARCHAR(50),
    fcreateorgid VARCHAR(50),
    fjoindate TIMESTAMP,
    fuseorgid VARCHAR(50),
    fname VARCHAR(100)
);

