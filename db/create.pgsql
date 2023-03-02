-- удаление схемы со всеми данными(если такая существует)
drop schema if exists eco_pasport;

-- создание схемы для данных
create schema if not exists eco_pasport;

-- создание таблицы заголовков(таблиц пасспортов)
create table eco_pasport."Header"(
	id int4 primary key,
	name varchar(80) not null,
	"header" json not null
);

-- создание таблицы регионов
create table eco_pasport."Region"(
	id int4 primary key,
	name varchar(80) not null,
	lat float4 not null,
	lng float4 not null,
	sort_region int4,
	inform_id int4,
	istown boolean

);

-- создание таблицы с доп информацией регионов
create table eco_pasport."Region_inform"(
	id int4 primary key,
	center varchar(80),
	data_create int4,
	population_size int4,
	area float4,
	gross_emissions varchar(80),
	withdrawn_water varchar(80),
	discharge_volume varchar(80),
	formed_waste varchar(80)
);

-- создание таблицы с основными данными пасспортов
create table eco_pasport."Row"(
	id int4 primary key,
	Value _varchar,
	Sort serial,
	Table_id int4 not null
);

-- создание линковой таблицы с информацией таблиц пасспортов
create table eco_pasport."Table"(
	id int4 primary key,
	"year" int4,
	header_id int4,
	region_id int4
);

-- Создание внешних ключей между таблицами
ALTER TABLE eco_pasport."Row" ADD CONSTRAINT row_fk FOREIGN KEY (table_id) REFERENCES eco_pasport."Table"(id);
ALTER TABLE eco_pasport."Table" ADD CONSTRAINT table_fk FOREIGN KEY (header_id) REFERENCES eco_pasport."Header"(id);
ALTER TABLE eco_pasport."Table" ADD CONSTRAINT table_fk_1 FOREIGN KEY (region_id) REFERENCES eco_pasport."Region"(id);
ALTER TABLE eco_pasport."Region" ADD CONSTRAINT region_fk FOREIGN KEY (inform_id) REFERENCES eco_pasport."Region_inform"(id);

-- Создание автоинкремента для таблиц
create sequence eco_pasport.region_id_seq
ALTER TABLE eco_pasport."Region" ALTER COLUMN id SET DEFAULT nextval('eco_pasport.region_id_seq'::regclass);
create sequence eco_pasport.header_id_seq
ALTER TABLE eco_pasport."Header" ALTER COLUMN id SET DEFAULT nextval('eco_pasport.header_id_seq'::regclass);
create sequence eco_pasport.row_id_seq
ALTER TABLE eco_pasport."Row" ALTER COLUMN id SET DEFAULT nextval('eco_pasport.row_id_seq'::regclass);
create sequence eco_pasport.table_id_seq
ALTER TABLE eco_pasport."Table" ALTER COLUMN id SET DEFAULT nextval('eco_pasport.table_id_seq'::regclass);
create sequence eco_pasport.regionInf_id_seq
ALTER TABLE eco_pasport."Region_inform" ALTER COLUMN id SET DEFAULT nextval('eco_pasport.regionInf_id_seq'::regclass);
