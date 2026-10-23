import random
import time

from locust import HttpUser, between, task


class Mumnk8sUser(HttpUser):
    wait_time = between(1, 3)

    @task
    def enviar_reporte(self):
        reporte = {
            "country_code": "ESP",
            "warplanes_in_air": random.randint(0, 50),
            "warships_in_water": random.randint(0, 30),
            "timestamp_unix": int(time.time()),
        }

        with self.client.post(
            "/reports",
            json=reporte,
            name="POST /reports",
            catch_response=True,
        ) as respuesta:
            if respuesta.status_code == 202:
                respuesta.success()
            else:
                respuesta.failure(
                    f"HTTP {respuesta.status_code}: {respuesta.text}"
                )
