package game

var LocationMap = map[uint8]string{
	0x00: "Littleroot Town",
	0x01: "Oldale Town",
	0x02: "Dewford Town",
	0x03: "Lavaridge Town",
	0x04: "Fallarbor Town",
	0x05: "Verdanturf Town",
	0x06: "Pacifidlog Town",
	0x07: "Petalburg City",
	0x08: "Slateport City",
	0x09: "Mauville City",
	0x0A: "Rustboro City",
	0x0B: "Fortree City",
	0x0C: "Lilycove City",
	0x0D: "Mossdeep City",
	0x0E: "Sootopolis City",
	0x0F: "Ever Grande City",
	0x10: "Route 101",
	0x11: "Route 102",
	0x12: "Route 103",
	0x13: "Route 104",
	0x14: "Route 105",
	0x15: "Route 106",
	0x16: "Route 107",
	0x17: "Route 108",
	0x18: "Route 109",
	0x19: "Route 110",
	0x1A: "Route 111",
	0x1B: "Route 112",
	0x1C: "Route 113",
	0x1D: "Route 114",
	0x1E: "Route 115",
	0x1F: "Route 116",
	0x20: "Route 117",
	0x21: "Route 118",
	0x22: "Route 119",
	0x23: "Route 120",
	0x24: "Route 121",
	0x25: "Route 122",
	0x26: "Route 123",
	0x27: "Route 124",
	0x28: "Route 125",
	0x29: "Route 126",
	0x2A: "Route 127",
	0x2B: "Route 128",
	0x2C: "Route 129",
	0x2D: "Route 130",
	0x2E: "Route 131",
	0x2F: "Route 132",
	0x30: "Route 133",
	0x31: "Route 134",
	0x32: "Underwater 124",
	0x33: "Underwater 125",
	0x34: "Underwater 126",
	0x35: "Underwater 127",
	0x36: "Underwater Sootopolis",
	0x37: "Granite Cave",
	0x38: "Mt Chimney",
	0x39: "Safari Zone",
	0x3A: "Battle Frontier",
	0x3B: "Petalburg Woods",
	0x3C: "Rusturf Tunnel",
	0x3D: "Abandoned Ship",
	0x3E: "New Mauville",
	0x3F: "Meteor Falls",
	0x40: "Meteor Falls2",
	0x41: "Mt Pyre",
	0x42: "Aqua Hideout Old",
	0x43: "Shoal Cave",
	0x44: "Seafloor Cavern",
	0x45: "Underwater 128",
	0x46: "Victory Road",
	0x47: "Mirage Island",
	0x48: "Cave Of Origin",
	0x49: "Southern Island",
	0x4A: "Fiery Path",
	0x4B: "Fiery Path2",
	0x4C: "Jagged Pass",
	0x4D: "Jagged Pass2",
	0x4E: "Sealed Chamber",
	0x4F: "Underwater Sealed Chamber",
	0x50: "Scorched Slab",
	0x51: "Island Cave",
	0x52: "Desert Ruins",
	0x53: "Ancient Tomb",
	0x54: "Inside Of Truck",
	0x55: "Sky Pillar",
	0x56: "Secret Base",
	0x57: "Dynamic",
	0x58: "Pallet Town",
	0x59: "Viridian City",
	0x5A: "Pewter City",
	0x5B: "Cerulean City",
	0x5C: "Lavender Town",
	0x5D: "Vermilion City",
	0x5E: "Celadon City",
	0x5F: "Fuchsia City",
	0x60: "Cinnabar Island",
	0x61: "Indigo Plateau",
	0x62: "Saffron City",
	0x63: "Route 4 Pokecenter",
	0x64: "Route 10 Pokecenter",
	0x65: "Route 1",
	0x66: "Route 2",
	0x67: "Route 3",
	0x68: "Route 4",
	0x69: "Route 5",
	0x6A: "Route 6",
	0x6B: "Route 7",
	0x6C: "Route 8",
	0x6D: "Route 9",
	0x6E: "Route 10",
	0x6F: "Route 11",
	0x70: "Route 12",
	0x71: "Route 13",
	0x72: "Route 14",
	0x73: "Route 15",
	0x74: "Route 16",
	0x75: "Route 17",
	0x76: "Route 18",
	0x77: "Route 19",
	0x78: "Route 20",
	0x79: "Route 21",
	0x7A: "Route 22",
	0x7B: "Route 23",
	0x7C: "Route 24",
	0x7D: "Route 25",
	0x7E: "Viridian Forest",
	0x7F: "Mt Moon",
	0x80: "SS Anne",
	0x81: "Underground Path",
	0x82: "Underground Path 2",
	0x83: "Digletts Cave",
	0x84: "Kanto Victory Road",
	0x85: "Rocket Hideout",
	0x86: "Silph Co",
	0x87: "Pokemon Mansion",
	0x88: "Kanto Safari Zone",
	0x89: "Pokemon League",
	0x8A: "Rock Tunnel",
	0x8B: "Seafoam Islands",
	0x8C: "Pokemon Tower",
	0x8D: "Cerulean Cave",
	0x8E: "Power Plant",
	0x8F: "One Island",
	0x90: "Two Island",
	0x91: "Three Island",
	0x92: "Four Island",
	0x93: "Five Island",
	0x94: "Seven Island",
	0x95: "Six Island",
	0x96: "Kindle Road",
	0x97: "Treasure Beach",
	0x98: "Cape Brink",
	0x99: "Bond Bridge",
	0x9A: "Three Isle Port",
	0x9B: "Sevii Isle 6",
	0x9C: "Sevii Isle 7",
	0x9D: "Sevii Isle 8",
	0x9E: "Sevii Isle 9",
	0x9F: "Resort Gorgeous",
	0xA0: "Water Labyrinth",
	0xA1: "Five Isle Meadow",
	0xA2: "Memorial Pillar",
	0xA3: "Outcast Island",
	0xA4: "Green Path",
	0xA5: "Water Path",
	0xA6: "Ruin Valley",
	0xA7: "Trainer Tower",
	0xA8: "Canyon Entrance",
	0xA9: "Sevault Canyon",
	0xAA: "Tanoby Ruins",
	0xAB: "Sevii Isle 22",
	0xAC: "Sevii Isle 23",
	0xAD: "Sevii Isle 24",
	0xAE: "Navel Rock",
	0xAF: "Mt Ember",
	0xB0: "Berry Forest",
	0xB1: "Icefall Cave",
	0xB2: "Rocket Warehouse",
	0xB3: "Trainer Tower 2",
	0xB4: "Dotted Hole",
	0xB5: "Lost Cave",
	0xB6: "Pattern Bush",
	0xB7: "Altering Cave",
	0xB8: "Tanoby Chambers",
	0xB9: "Three Isle Path",
	0xBA: "Tanoby Key",
	0xBB: "Birth Island",
	0xBC: "Monean Chamber",
	0xBD: "Liptoo Chamber",
	0xBE: "Weepth Chamber",
	0xBF: "Dilford Chamber",
	0xC0: "Scufib Chamber",
	0xC1: "Rixy Chamber",
	0xC2: "Viapois Chamber",
	0xC3: "Ember Spa",
	0xC4: "Special Area",
	0xC5: "None",
	0xC6: "Count",
	0xFD: "(Egg)",
	0xFE: "(Trade)",
	0xFF: "(Fateful)",
}
