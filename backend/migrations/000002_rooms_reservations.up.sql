CREATE TABLE room_types (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    image TEXT NOT NULL,
    max_adults INT NOT NULL CHECK (max_adults >= 1),
    max_children INT NOT NULL CHECK (max_children >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE room_prices (
    room_type_id TEXT NOT NULL REFERENCES room_types(id) ON DELETE CASCADE,
    occupancy INT NOT NULL CHECK (occupancy BETWEEN 1 AND 4),
    price_from NUMERIC(10,2) NOT NULL,
    price_to NUMERIC(10,2) NOT NULL,
    PRIMARY KEY (room_type_id, occupancy)
);

CREATE TABLE reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NULL REFERENCES users(id) ON DELETE SET NULL,
    guest_name TEXT NOT NULL,
    guest_phone TEXT NOT NULL,
    guest_email TEXT NOT NULL,
    check_in DATE NOT NULL,
    check_out DATE NOT NULL,
    nights INT NOT NULL CHECK (nights >= 1),
    total NUMERIC(12,2) NOT NULL,
    status TEXT NOT NULL DEFAULT 'confirmed',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT reservations_dates_chk CHECK (check_out > check_in)
);

CREATE TABLE reservation_rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reservation_id UUID NOT NULL REFERENCES reservations(id) ON DELETE CASCADE,
    room_type_id TEXT NOT NULL REFERENCES room_types(id),
    adults INT NOT NULL CHECK (adults >= 1),
    children INT NOT NULL CHECK (children >= 0),
    nightly_price NUMERIC(10,2) NOT NULL,
    line_total NUMERIC(12,2) NOT NULL
);

CREATE INDEX reservations_user_id_idx ON reservations(user_id);
CREATE INDEX reservation_rooms_reservation_id_idx ON reservation_rooms(reservation_id);

INSERT INTO room_types (id, name, description, image, max_adults, max_children) VALUES
('express-duplo', 'Express Duplo', 'Modelo com 02 camas de solteiro, frigobar, TV e ar-condicionado.', '/assets/img/twin-bed.jpg', 2, 2),
('express-triplo', 'Express Triplo', 'Modelo com 03 camas de solteiro, frigobar, TV e ar-condicionado.', '/assets/img/triplo.jpg', 3, 2),
('luxo-casal', 'Luxo Casal', 'Modelo com 01 cama de casal, frigobar, TV e ar-condicionado.', '/assets/img/double-bed.jpg', 2, 2),
('luxo-familia', 'Luxo Família', 'Modelo com 01 cama de casal, 01 cama de solteiro, frigobar, TV e ar-condicionado.', '/assets/img/double_bed_and_single_bed.jpg', 3, 2),
('flat', 'Flat', 'Modelo com 01 cama de casal, 02 camas de solteiro, mini-cozinha, sofá, TV e ar-condicionado.', '/assets/img/flat.jpg', 4, 4),
('master', 'Master', 'Modelo com 01 cama de casal, sofá, TV, ar-condicionado e hidromassagem.', '/assets/img/master.jpg', 2, 2);

INSERT INTO room_prices (room_type_id, occupancy, price_from, price_to) VALUES
('express-duplo', 1, 200, 180),
('express-duplo', 2, 269, 249),
('express-triplo', 1, 250, 210),
('express-triplo', 2, 269, 249),
('express-triplo', 3, 349, 329),
('luxo-casal', 1, 249, 229),
('luxo-casal', 2, 299, 279),
('luxo-familia', 1, 249, 229),
('luxo-familia', 2, 299, 279),
('luxo-familia', 3, 369, 349),
('flat', 1, 330, 310),
('flat', 2, 410, 390),
('flat', 3, 530, 510),
('flat', 4, 600, 580),
('master', 1, 400, 380),
('master', 2, 520, 500);
