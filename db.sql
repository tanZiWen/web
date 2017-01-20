CREATE TABLE if not exists client
(
  id text NOT NULL,
  secret text NOT NULL,
  extra text NOT NULL,
  redirect_uri text NOT NULL,
  CONSTRAINT client_pkey PRIMARY KEY (id)
);

CREATE TABLE if not exists authorize
(
  code text NOT NULL,
  client text NOT NULL,
  expires_in integer NOT NULL,
  scope text NOT NULL,
  redirect_uri text NOT NULL,
  state text NOT NULL,
  extra text NOT NULL,
  created_at timestamp with time zone,
  CONSTRAINT authorize_pkey PRIMARY KEY (code)
);

CREATE TABLE if not exists access
(
  access_token text NOT NULL,
  client text NOT NULL,
  authorize text NOT NULL,
  previous text NOT NULL,
  refresh_token text NOT NULL,
  expires_in integer NOT NULL,
  scope text NOT NULL,
  redirect_uri text NOT NULL,
  extra text NOT NULL,
  created_at timestamp with time zone,
  CONSTRAINT access_pkey PRIMARY KEY (access_token)
);

CREATE TABLE if not exists refresh
(
  token text NOT NULL,
  access text NOT NULL,
  CONSTRAINT refresh_pkey PRIMARY KEY (token)
);

CREATE TABLE upm_user
(
  user_id bigint NOT NULL, -- 用户id
  workgroup_id bigint, -- 团队id
  workgroup_role_codes character varying[], -- 团队角色编码
  status character varying(30), -- 状态
  user_name character varying(30), -- 用户名
  password character varying(1000), -- 密码
  real_name character varying(30), -- 姓名
  sex character varying(10), -- 性别
  email character varying(30), -- 邮箱
  position character varying(30), -- 职位
  employee_id bigint NOT NULL, -- 员工id
  employee_code character varying(30), -- 员工编号
  work_no character varying(30), -- 工作编号
  ext_no character varying(30), -- 分机号
  org_code character varying(30), -- 机构编号
  role_codes character varying[], -- 角色编码
  online boolean, -- 是否在线
  c_user_id bigint, -- 创建人id
  c_user_name character varying(30), -- 创建人用户名
  c_real_name character varying(30), -- 创建人真实名
  c_time time(6) without time zone, -- 创建时间
  u_user_id bigint, -- 最后修改人id
  u_user_name character varying(30), -- 最后修改人用户名
  u_real_name character varying(30), -- 最后修改人真实名
  u_time time(6) without time zone, -- 最后修改时间
  CONSTRAINT upm_user_pkey PRIMARY KEY (user_id)
);
COMMENT ON COLUMN upm_user.user_id IS '用户id';
COMMENT ON COLUMN upm_user.workgroup_id IS '团队id';
COMMENT ON COLUMN upm_user.workgroup_role_codes IS '团队角色编码';
COMMENT ON COLUMN upm_user.status IS '状态';
COMMENT ON COLUMN upm_user.user_name IS '用户名';
COMMENT ON COLUMN upm_user.password IS '密码';
COMMENT ON COLUMN upm_user.real_name IS '姓名';
COMMENT ON COLUMN upm_user.sex IS '性别';
COMMENT ON COLUMN upm_user.email IS '邮箱';
COMMENT ON COLUMN upm_user."position" IS '职位';
COMMENT ON COLUMN upm_user.employee_id IS '员工id';
COMMENT ON COLUMN upm_user.employee_code IS '员工编号';
COMMENT ON COLUMN upm_user.work_no IS '工作编号';
COMMENT ON COLUMN upm_user.ext_no IS '分机号';
COMMENT ON COLUMN upm_user.org_code IS '机构编号';
COMMENT ON COLUMN upm_user.role_codes IS '角色编码';
COMMENT ON COLUMN upm_user.online IS '是否在线';
COMMENT ON COLUMN upm_user.c_user_id IS '创建人id';
COMMENT ON COLUMN upm_user.c_user_name IS '创建人用户名';
COMMENT ON COLUMN upm_user.c_real_name IS '创建人真实名';
COMMENT ON COLUMN upm_user.c_time IS '创建时间';
COMMENT ON COLUMN upm_user.u_user_id IS '最后修改人id';
COMMENT ON COLUMN upm_user.u_user_name IS '最后修改人用户名';
COMMENT ON COLUMN upm_user.u_real_name IS '最后修改人真实名';
COMMENT ON COLUMN upm_user.u_time IS '最后修改时间';
