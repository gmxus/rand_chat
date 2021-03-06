package main

import (
	"github.com/droptheplot/rand_chat/env"
)

func main() {
	env.DB.Exec(`
		DO $$
		DECLARE
			room_id integer;
			i integer = 0;
		BEGIN
			INSERT INTO rooms (owner_id, guest_id) VALUES (random() * 100, random() * 100) RETURNING id INTO room_id;

			WHILE i < 10
			LOOP
				i := i + 1;

				INSERT INTO messages (room_id, created_at) VALUES
					(room_id, now() - (i * (random() * 100) || ' minutes')::interval);
			END LOOP;
		END $$;
	`)
}
