package database

import (
	"log"
	"net/http"
    "github.com/clairerhoda/Database_CPSC332/department"
    "github.com/clairerhoda/Database_CPSC332/building"
    "github.com/clairerhoda/Database_CPSC332/office"
    "github.com/clairerhoda/Database_CPSC332/user"
    "github.com/clairerhoda/Database_CPSC332/user_proxy"
    "github.com/clairerhoda/Database_CPSC332/sectional_room"
    "github.com/clairerhoda/Database_CPSC332/room"
    "github.com/clairerhoda/Database_CPSC332/equipment"
    "github.com/clairerhoda/Database_CPSC332/reservation"
    "github.com/clairerhoda/Database_CPSC332/reservation_attendee"
    "github.com/clairerhoda/Database_CPSC332/reserved_equipment"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func Execute() {
	r := mux.NewRouter()
	department.DepartmentExecute()
	building.BuildingExecute(r)
	office.OfficeExecute(r)
	user.UserExecute(r)
	user_proxy.UserProxyExecute(r)
	sectional_room.SectionalRoomExecute()
	room.RoomExecute(r)
	equipment.EquipmentExecute(r)
	reservation.ReservationExecute(r)
	reservation_attendee.ReservationAttendeeExecute(r)
	reserved_equipment.ReservedEquipmentExecute(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}