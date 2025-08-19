package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const (
	// getDeviceByID is a SQL string to select a device by its ID. It returns the device's ID, user_id, device_name, and created_at.
	getDeviceByID string = `SELECT id, user_id, device_name, created_at FROM devices WHERE id = $1;`
	// getDevicesByUserID is a SQL string to select all devices for a user by their user ID. It returns the devices' ID, user_id, device_name, and created_at.
	getDevicesByUserID string = `SELECT id, user_id, device_name, created_at FROM devices WHERE user_id = $1;`
	// saveDevice is a SQL string to insert a new device into the database. It returns the newly created device's ID.
	saveDevice string = `INSERT INTO devices (user_id, device_name) VALUES ($1, $2) RETURNING id;`
)

type Device struct {
	ID        uuid.UUID
	User      *User
	Name      string
	CreatedAt string
}

func (d *Device) unmarshalRow(row pgx.Row) error {
	return row.Scan(&d.ID, &d.User.ID, &d.Name, &d.CreatedAt)
}

func (db *DB) GetDeviceByID(id uuid.UUID) (*Device, error) {
	var d Device
	row := db.conn.QueryRow(context.Background(), getDeviceByID, id)
	if err := d.unmarshalRow(row); err != nil {
		return nil, err
	}
	return &d, nil
}

func (db *DB) GetDevicesByUserID(userID uuid.UUID) ([]*Device, error) {
	rows, err := db.conn.Query(context.Background(), getDevicesByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*Device
	for rows.Next() {
		var d Device
		if err := d.unmarshalRow(rows); err != nil {
			return nil, err
		}
		devices = append(devices, &d)
	}

	return devices, nil
}

func (db *DB) InsertDevice(userID uuid.UUID, deviceName string) (uuid.UUID, error) {
	var id uuid.UUID
	row := db.conn.QueryRow(context.Background(), saveDevice, userID, deviceName)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
