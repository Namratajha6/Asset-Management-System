CREATE TYPE asset_type AS ENUM ('laptop', 'mouse', 'hard_disk', 'pendrive', 'mobile', 'sim', 'accessories');
CREATE TYPE asset_status AS ENUM ('available', 'assigned', 'waiting_for_repair', 'service', 'damaged');
CREATE TYPE employee_type AS ENUM ('employee', 'intern', 'freelance');
CREATE TYPE employee_role AS ENUM ('admin', 'asset_manager', 'employee_manager', 'employee');
CREATE TYPE asset_owned_by AS ENUM ('remote_state', 'client');
CREATE TYPE asset_history_status AS ENUM ('assigned', 'retrieved');

CREATE TABLE IF NOT EXISTS employees (
                                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                         name TEXT NOT NULL,
                                         email TEXT NOT NULL UNIQUE,
                                         phone_no TEXT,
                                         type employee_type NOT NULL,
                                         role employee_role NOT NULL,
                                         created_at TIMESTAMPTZ DEFAULT NOW(),
                                         archived_at TIMESTAMPTZ DEFAULT NULL,
                                         created_by UUID REFERENCES employees(id),
                                         updated_by UUID REFERENCES employees(id)
);

CREATE TABLE IF NOT EXISTS assets (
                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                      brand TEXT NOT NULL,
                                      model TEXT NOT NULL,
                                      serial_no TEXT NOT NULL UNIQUE,
                                      asset_type asset_type NOT NULL,
                                      asset_status asset_status DEFAULT 'available',
                                      owned_by asset_owned_by NOT NULL DEFAULT 'remote_state',
                                      purchased_date DATE,
                                      warranty_start_date DATE,
                                      warranty_end_date DATE,
                                      created_at TIMESTAMPTZ DEFAULT NOW(),
                                      archived_at TIMESTAMPTZ DEFAULT NULL,
                                      created_by UUID REFERENCES employees(id) NOT NULL,
                                      updated_by UUID REFERENCES employees(id)
);

CREATE TABLE IF NOT EXISTS laptops (
                                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                       asset_id UUID UNIQUE REFERENCES assets(id) ,
                                       os TEXT NOT NULL,
                                       ram TEXT NOT NULL,
                                       storage TEXT NOT NULL,
                                       processor TEXT NOT NULL,
                                       created_at TIMESTAMPTZ DEFAULT NOW(),
                                       archived_at TIMESTAMPTZ DEFAULT NULL,
                                       created_by UUID REFERENCES employees(id) NOT NULL,
                                       updated_by UUID REFERENCES employees(id)
);

CREATE TABLE IF NOT EXISTS mouse (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     asset_id UUID UNIQUE REFERENCES assets(id) ,
                                     connectivity_type TEXT NOT NULL,
                                     created_at TIMESTAMPTZ DEFAULT NOW(),
                                     archived_at TIMESTAMPTZ DEFAULT NULL,
                                     created_by UUID REFERENCES employees(id) NOT NULL,
                                     updated_by UUID REFERENCES employees(id)
);

CREATE TABLE IF NOT EXISTS hard_disks (
                                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                         asset_id UUID UNIQUE REFERENCES assets(id) ,
                                         storage_capacity TEXT NOT NULL,
                                         created_at TIMESTAMPTZ DEFAULT NOW(),
                                         archived_at TIMESTAMPTZ DEFAULT NULL,
                                         created_by UUID REFERENCES employees(id) NOT NULL,
                                         updated_by UUID REFERENCES employees(id)
);

CREATE TABLE IF NOT EXISTS pendrives (
                                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                         asset_id UUID UNIQUE REFERENCES assets(id) ,
                                         storage_capacity TEXT NOT NULL,
                                         created_at TIMESTAMPTZ DEFAULT NOW(),
                                         archived_at TIMESTAMPTZ DEFAULT NULL,
                                         created_by UUID REFERENCES employees(id) NOT NULL,
                                         updated_by UUID REFERENCES employees(id)
);

CREATE TABLE IF NOT EXISTS mobiles (
                                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                       asset_id UUID UNIQUE REFERENCES assets(id) ,
                                       imei1 TEXT NOT NULL,
                                       imei2 TEXT NOT NULL,
                                       os TEXT NOT NULL,
                                       ram TEXT NOT NULL,
                                       storage TEXT NOT NULL,
                                       created_at TIMESTAMPTZ DEFAULT NOW(),
                                       archived_at TIMESTAMPTZ DEFAULT NULL,
                                       created_by UUID REFERENCES employees(id) NOT NULL,
                                       updated_by UUID REFERENCES employees(id)
);

CREATE TABLE IF NOT EXISTS sims (
                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                    asset_id UUID UNIQUE REFERENCES assets(id) ,
                                    mobile_number TEXT NOT NULL UNIQUE,
                                    network_provider TEXT NOT NULL,
                                    created_at TIMESTAMPTZ DEFAULT NOW(),
                                    archived_at TIMESTAMPTZ DEFAULT NULL,
                                    created_by UUID REFERENCES employees(id) NOT NULL,
                                    updated_by UUID REFERENCES employees(id)
);

CREATE TABLE IF NOT EXISTS asset_employee_history (
                                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                                      employee_id UUID REFERENCES employees(id) ,
                                                      asset_id UUID REFERENCES assets(id) ,
                                                      assigned_date DATE,
                                                      return_date DATE,
                                                      status asset_history_status,
                                                      performed_at TIMESTAMPTZ DEFAULT NOW(),
                                                      performed_by UUID REFERENCES employees(id) NOT NULL,
                                                      updated_by UUID REFERENCES employees(id)
);

CREATE TABLE IF NOT EXISTS asset_history (
                                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                             asset_id UUID REFERENCES assets(id) ,
                                             old_status asset_status,
                                             new_status asset_status NOT NULL,
                                             employee_id UUID REFERENCES employees(id),
                                             performed_by UUID REFERENCES employees(id),
                                             performed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_employees_email ON employees(email) WHERE archived_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_asset_emp_history_asset_id ON asset_employee_history(asset_id);
CREATE INDEX IF NOT EXISTS idx_asset_emp_history_employee_id ON asset_employee_history(employee_id);

CREATE INDEX IF NOT EXISTS idx_asset_employee_performed_at
    ON asset_employee_history(asset_id, performed_at DESC);

