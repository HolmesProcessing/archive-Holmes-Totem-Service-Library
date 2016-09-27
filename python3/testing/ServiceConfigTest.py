import unittest
from python3.services.configuration import ParseConfig


exampleConfiguration = """
{
    "Port": 5666,
    "IP": "53.178.35.27",
    "Limit": 5000,
    "ExtraSettings": {
        "keyA": "A",
        "keyB": "B",
        "keyC": "C"
    }
}
"""


class TestServiceConfig(unittest.TestCase):

    def test_all(self):
        config = {
            "port": 8016,
            "ip": "0.0.0.0",
            "ip2": "0.0.0.0",
            "limit": 1000,
            "Extrasettings": {
                "KeyA": "--empty--",
                "KeyB": "--empty--",
                "KeyC": "--empty--",
                "KeyD": "--empty--",
            }
        }
        cfg = ParseConfig(config, data=exampleConfiguration)

        self.assertEqual(config["port"],  5666)
        self.assertEqual(config["ip"],    "53.178.35.27")
        self.assertEqual(config["ip2"],   "0.0.0.0")
        self.assertEqual(config["limit"], 5000)
        self.assertEqual(config["Extrasettings"]["KeyA"], "A")
        self.assertEqual(config["Extrasettings"]["KeyB"], "B")
        self.assertEqual(config["Extrasettings"]["KeyC"], "C")
        self.assertEqual(config["Extrasettings"]["KeyD"], "--empty--")

        self.assertEqual(cfg.PORT,  5666)
        self.assertEqual(cfg.IP2,   "0.0.0.0")
        self.assertEqual(cfg.LimiT, 5000)
        self.assertEqual(cfg.extrasettingS.keya, "A")
        self.assertEqual(cfg.extrasettinGs.KEYB, "B")
        self.assertEqual(cfg.extrasettiNgs.keYC, "C")


if __name__ == '__main__':
    unittest.main()
