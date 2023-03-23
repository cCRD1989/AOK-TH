package khanscr

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var dataFolderPath string
var itemNames map[int16]string
var itemNamesMsg []string
var itemIDs map[string]int16
var bsvItems map[int16]BsvItem
var items = make(map[int16]Item)

var (
	regNumFind = regexp.MustCompile(`"[^"]+"|([-0-9]+)`)
	regNewLine = regexp.MustCompile(`[\r\n]+`)
	regTabs    = regexp.MustCompile(`[\t]+`)
)

// BsvItem structure
// # T T i W W B W5 B2 W B W6 D3 B2 W s B W4 W5 D i W W i i2 i i
// FORMAT_ITEM_INFO = "=ihhhhhhhhbbhhhhhhhiiiibbhhhhhhhhhhhiiihhiiiii"
// i = int(4)
// h = short(2)
// b = signed char(1)
type BsvItem struct {
	Unknown1              int32
	ItemType              int16
	Wield                 int16
	Description           int16
	Class                 int16
	DamageMinMelee        int16
	DamageMaxMelee        int16
	DamageMinMagic        int16
	DamageMaxMagic        int16
	Range                 byte
	Durability            byte
	AttackSpeed           int16
	Unknown13             int16
	Strength              int16
	Dexterity             int16
	Wisdom                int16
	Charisma              int16
	LevelRequirement      int16
	LevelRequirementRange int32
	Unknown20             int32
	ItemSellPrice         int32
	ItemBuyPrice          int32
	ColCell               byte
	RowCell               byte
	ItemLimitation        int16
	Unknown26             int16
	Unknown27             int16
	ItemName              int16
	Enchant               int16
	EffectMessage         int16
	Unknown31             int16
	ItemModel             int16
	Unknown33             int16
	EquipPosition         int16
	ItemIcon              int16
	Unknown3637           [2]int32
	Effect                int32
	Unknown39             int16
	Unknown40             int16
	RecipeA               int32
	RecipeB               int32
	RecipeC               int32
	UpGrade               int32
	Unknown45             int32
}

// Item object
type Item struct {
	BsvItem
	ID   int16
	Name string
}

// Init initializes the khan data
func Init() {
	// root path
	dataFolderPath = "./server/khandata/GameServer/Common_Data/"

	// initialize maps
	itemNames = make(map[int16]string) // 0:broad sword
	itemIDs = make(map[string]int16)   // broad sword:0
	bsvItems = make(map[int16]BsvItem)

	// load items --------------------------------------------
	// load ItemInfo_Name
	var file = readFile(cmnPath("ItemInfo_Name.txt"))

	readDataRaw(file, func(line string) {
		var itemInfo = strings.Split(line, "\t")
		var itemID, _ = strconv.Atoi(itemInfo[0])
		itemNames[int16(itemID)] = itemInfo[1]

	})
	for k, v := range itemNames {
		itemIDs[v] = k
	}

	// load iteminfo.BSV
	var bsv = readFile(cmnPath("iteminfo.bsv"))
	defer bsv.Close()

	bsv.Seek(32, os.SEEK_SET)
	var fileStart = int64(readUInt32(bsv))
	bsv.Seek(44, os.SEEK_SET)
	var numRows = int(readUInt32(bsv))
	bsv.Seek(fileStart, os.SEEK_SET)

	var bsvItem BsvItem

	for i := 0; i < numRows; i++ {
		var data = readNextBytes(bsv, 112)
		var buffer = bytes.NewBuffer(data)
		var err = binary.Read(buffer, binary.LittleEndian, &bsvItem)

		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		bsvItems[int16(i)] = bsvItem
		items[int16(i)] = Item{
			BsvItem: bsvItem,
			ID:      int16(i),
			Name:    itemNames[int16(i)],
		}
	}

	// itemNames
	var itemNamesBin = readFile(clientPath("itemname.bin"))
	readData(itemNamesBin, func(i int, l string) {
		line := strings.Split(regTabs.ReplaceAllString(l, "\t"), "\t")
		itemName := line[1]
		itemNamesMsg = append(itemNamesMsg, itemName)
	})
	// End load items --------------------------------------------
}

// GetAllItemNames returns all items
func GetAllItemNames() map[int16]string {
	return itemNames
}

// GetItemName get the name of the item
func GetItemName(itemID int16) string {
	return itemNames[itemID]
}

// GetItemIDByName get the id by name
func GetItemIDByName(name string) int16 {
	name = strings.ToLower(strings.TrimSpace(name))
	if itemID, ok := itemIDs[name]; ok {
		return itemID
	}
	return -1
}

// GetAllItemNamesMsg returns name from the .bin msg file
func GetAllItemNamesMsg() []string {
	return itemNamesMsg
}

// GetAllItems returns all item data
func GetAllItems() map[int16]BsvItem {
	return bsvItems
}

// GetItemByID returns a BsvItem structs
func GetItemByID(itemID int16) BsvItem {
	return bsvItems[itemID]
}

// GetAllItemsWithName returns all item data
func GetAllItemsWithName() map[int16]Item {
	return items
}

/////////////////////////////////////////////////////////////////

func cmnPath(file string) string {
	var path, _ = filepath.Abs(fmt.Sprintf("%s/%s", dataFolderPath, file))
	return path
}

func clientPath(file string) string {
	var path, _ = filepath.Abs(fmt.Sprintf("%s/client/%s", dataFolderPath, file))
	return path
}

func readFile(filepath string) *os.File {
	var file, err = os.Open(filepath)

	if err != nil {
		log.Fatal("func readFile ", err)
	}

	return file
}

func readDataRaw(bin *os.File, fn func(line string)) {
	defer bin.Close()

	var scanner = bufio.NewScanner(bin)
	for scanner.Scan() {
		var line = scanner.Text()
		// ignore comments
		if string(line[0]) != "#" {
			fn(line)
		}
	}
}
func readUInt32(file *os.File) uint32 {
	bytes := readNextBytes(file, 4)
	return binary.LittleEndian.Uint32(bytes)
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

// Bin Client
func readData(bin *os.File, fn func(i int, line string)) {
	data, err := ioutil.ReadAll(bin)
	if err != nil {
		log.Fatal(err)
	}
	defer bin.Close()

	for i := 0; i < len(data); i++ {
		if string(data[i]) != "\r" && string(data[i]) != "\n" {
			data[i] = data[i] ^ 0x14
		}
	}

	var lines = regNewLine.Split(string(data), -1)
	lines = deleteEmpty(lines)

	for idx, line := range lines {
		// ignore comments
		if string(line[0]) != "#" {
			fn(idx, line)
		}
	}
}
func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
