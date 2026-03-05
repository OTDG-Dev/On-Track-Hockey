import requests
import sys
import logging
from dataclasses import dataclass

logging.basicConfig(
    level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s"
)
URL = "http://localhost:3000/v1"


def test_conn():
    try:
        r = requests.get(URL)
    except Exception as e:
        print(e)
        sys.exit(1)


@dataclass
class Req:
    path: str
    payload: dict

    def __post_init__(self):
        self.path = URL + self.path


def post_request(path: str, payload: dict) -> bool:
    try:
        r = requests.post(url=path, json=payload)
        if r.status_code < 300:
            logging.info(f"created {path.split('/')[-1]}")
            return False
        logging.error(f"url: {r.url} status: {r.status_code} {r.text}")
        return True
    except Exception as e:
        print(f"path: {path}, payload {payload}, {e}")
        return True


def main():
    test_conn()

    request_list = [
        Req("/leagues", {"name": "league"}),
        Req("/divisions", {"league_id": 1, "name": "Team"}),
        Req(
            "/teams",
            {
                "full_name": "New York Rangers",
                "short_name": "NYR",
                "division_id": 1,
                "is_active": True,
            },
        ),
        Req(
            "/teams",
            {
                "full_name": "San Jose Sharks",
                "short_name": "SJS",
                "division_id": 1,
                "is_active": True,
            },
        ),
        Req(
            "/players",
            {
                "first_name": "Igor",
                "last_name": "Shesterkin",
                "sweater_number": 31,
                "position": "G",
                "birth_date": "1995-12-30",
                "birth_country": "RUS",
                "shoots_catches": "L",
                "current_team_id": 1,
                "is_active": True,
            },
        ),
        Req(
            "/players",
            {
                "first_name": "Macklin",
                "last_name": "Celebrini",
                "sweater_number": 71,
                "position": "C",
                "birth_date": "2006-05-12",
                "birth_country": "CAN",
                "shoots_catches": "L",
                "current_team_id": 2,
                "is_active": True,
            },
        ),
    ]

    for request in request_list:
        err = post_request(request.path, request.payload)


if __name__ == "__main__":
    main()
