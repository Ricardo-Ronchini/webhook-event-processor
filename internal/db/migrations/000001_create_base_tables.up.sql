CREATE TABLE inventory (
  inventory_id VARCHAR(30) PRIMARY KEY,
  product_id   VARCHAR(30) UNIQUE NOT NULL,
  sku          VARCHAR(30) NOT NULL,
  quantity     INTEGER NOT NULL DEFAULT 0,
  warehouse    VARCHAR(100),
  created_at   TIMESTAMP DEFAULT NOW(),
  updated_at   TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_inventory_product_id ON inventory (product_id);
CREATE INDEX idx_inventory_sku ON inventory (sku);

CREATE TABLE inventory_tracks (
  track_id     VARCHAR(30) PRIMARY KEY,
  inventory_id VARCHAR(30),
  product_id   VARCHAR(30) NOT NULL,
  event_type   VARCHAR(50) NOT NULL,
  quantity     INTEGER NOT NULL,
  payload      JSONB NOT NULL,
  created_at   TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_inventory_tracks_product_id ON inventory_tracks (product_id);
CREATE INDEX idx_inventory_tracks_created_at ON inventory_tracks (created_at);