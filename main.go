package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

// Game Plan:
// 1. Make the player appear
// 2. Player movement
// 3. Make sure player doesn't go off-screen
// 4. Zombies appear dynamically
// 5. Make zombies move
// 6. Fire bullets & make them move
// 7. Collision #1 - zombie & wall
// 8. Collision #2 - bullet & wall
// 9. Collision #3 - player is zombie
// 10. Collision #4 - bullet & zombie

const GameFrameWidth = 120
const GameFrameHeight = 25
const GameFrameSymbol = '🍃'
const GameFrameSymbol2 = '🍂'

type Point struct {
	row, col int
	symbol   rune
}

type GameObject struct {
	points         []*Point
	velRow, velCol int
	health         int
}

var screen tcell.Screen
var score int
var isGameOver bool
var isGamePaused bool

var player1 *GameObject
var player2 *GameObject

var zombies []*GameObject
var zombies2 []*GameObject
var zombies3 []*GameObject

var bullets []*GameObject

func main() {
	rand.Seed(time.Now().UnixNano())
	InitScreen()
	InitGameState()
	// Yourwin(score)

	inputChan := InitUserInput()

	for !isGameOver {
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()

		time.Sleep(100 * time.Millisecond)
	}

	width, height := screen.Size()
	PrintStringCentered(height/2-15, width/2-30, "🌵  Plants 🌻")
	PrintStringCentered(height/2-15, width/2, "🛠  V/S ⚔")
	PrintStringCentered(height/2-15, width/2+30, "👽  Zombies ☠")

	if score >= 200 {
		PrintStringCentered(height/2, width/2, "Plants win!!!")
		PrintStringCentered(height/2+1, width/2, fmt.Sprintf("You killed %d Zombies ", score))
		PrintStringCentered(height/2-16, width/2-53, "Plants are")
		PrintStringCentered(height/2-15, width/2-50, "the winners 🌹")
		PrintStringCentered(height/2+2, width/2, "🌞")
		PrintStringCentered(height/2+2, width/2+2, "🌞")
		PrintStringCentered(height/2+2, width/2+4, "🌞")
		PrintStringCentered(height/2, width/2+65, "❌")
		PrintStringCentered(height/2+5, width/2+65, "❌")
		PrintStringCentered(height/2-5, width/2+65, "❌")
		PrintStringCentered(height/2, width/2-67, "✅")
		PrintStringCentered(height/2+5, width/2-67, "✅")
		PrintStringCentered(height/2-5, width/2-67, "✅")
	} else {
		PrintStringCentered(height/2-15, width/2-50, "🥀")
		PrintStringCentered(height/2, width/2, "Game Over!!!")
		PrintStringCentered(height/2+1, width/2, fmt.Sprintf("You killed %d Zombies ", score))
		PrintStringCentered(height/2-16, width/2+45, "Zombies are")
		PrintStringCentered(height/2-15, width/2+45, "the winners")
		PrintStringCentered(height/2-15, width/2+60, "🏴‍☠️ 🏴")
		PrintStringCentered(height/2+2, width/2, "🌚")
		PrintStringCentered(height/2+2, width/2+2, "🌚")
		PrintStringCentered(height/2+2, width/2+4, "🌚")
		PrintStringCentered(height/2, width/2-67, "❌")
		PrintStringCentered(height/2+5, width/2-67, "❌")
		PrintStringCentered(height/2-5, width/2-67, "❌")
		PrintStringCentered(height/2, width/2+65, "✅")
		PrintStringCentered(height/2+5, width/2+65, "✅")
		PrintStringCentered(height/2-5, width/2+65, "✅")
	}

	if score <= 100 {
		PrintStringCentered(height/2-18, width/2-45, "😟")
		PrintStringCentered(height/2-18, width/2-42, "😢")
		PrintStringCentered(height/2-18, width/2-39, "🥵")
		PrintStringCentered(height/2-17, width/2-40, fmt.Sprintf("You killed %d Zombies, that's too few!!!", score))
	} else if score > 100 && score < 200 {
		PrintStringCentered(height/2-18, width/2-45, "🙂")
		PrintStringCentered(height/2-18, width/2-42, "🙂")
		PrintStringCentered(height/2-18, width/2-39, "🙂")
		PrintStringCentered(height/2-17, width/2-40, fmt.Sprintf("You killed %d Zombies, that's good!!!", score))
	} else if score >= 200 {
		PrintStringCentered(height/2-18, width/2-45, "😃")
		PrintStringCentered(height/2-18, width/2-42, "😃")
		PrintStringCentered(height/2-18, width/2-39, "😃")
		PrintStringCentered(height/2-17, width/2-40, fmt.Sprintf("You killed %d Zombies, that's great!!!", score))
	}
	screen.Show()
	time.Sleep(10 * time.Second)
	screen.Fini()
}

// func Yourwin(score int) bool {
// 	for score > 10 {
// 		isGameOver = false
// 	}
// 	return true
// }

func InitGameState() {
	player1 = &GameObject{
		points: []*Point{
			{row: 5, col: 2, symbol: '🌴'},
			{row: 6, col: 3, symbol: '█'},
			{row: 7, col: 3, symbol: '▀'},
		},
	}

	player2 = &GameObject{
		points: []*Point{
			{row: 10, col: 2, symbol: '🌾'},
			{row: 11, col: 3, symbol: '█'},
			{row: 12, col: 3, symbol: '▀'},
		},
	}
}

func InitScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defStyle := tcell.StyleDefault.
		Background(tcell.ColorGray).
		Foreground(tcell.ColorGoldenrod)
	screen.SetStyle(defStyle)
}

func InitUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()
	return inputChan
}

func ReadInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}
	return key
}

func HandleUserInput(key string) {
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
	} else if key == "Enter" {
		SplawnBullet(player1.points[0].row+1, player1.points[0].col+3)
	} else if key == "Rune[r]" {
		SplawnBullet2(player2.points[0].row, player2.points[0].col+3)
	} else if key == "Up" && !IsObjectOutOfBounds(player1, -1, 0) {
		NovePlayer(player1, -1, 0)
	} else if key == "Down" && !IsObjectOutOfBounds(player1, 1, 0) {
		NovePlayer(player1, 1, 0)
	} else if key == "Left" && !IsObjectOutOfBounds(player1, 0, -1) {
		NovePlayer(player1, 0, -1)
	} else if key == "Right" && !IsObjectOutOfBounds(player1, 0, 1) {
		NovePlayer(player1, 0, 1)
	}

	if key == "Rune[w]" && !IsObjectOutOfBounds(player2, -1, 0) {
		NovePlayer(player2, -1, 0)
	} else if key == "Rune[s]" && !IsObjectOutOfBounds(player2, 1, 0) {
		NovePlayer(player2, 1, 0)
	} else if key == "Rune[a]" && !IsObjectOutOfBounds(player2, 0, -1) {
		NovePlayer(player2, 0, -1)
	} else if key == "Rune[d]" && !IsObjectOutOfBounds(player2, 0, 1) {
		NovePlayer(player2, 0, 1)
	}
}

func NovePlayer(player *GameObject, velRow, velCol int) {
	for i := range player.points {
		player.points[i].row += velRow
		player.points[i].col += velCol
	}
}

func SplawnBullet(row, col int) {
	bullets = append(bullets, &GameObject{
		points: []*Point{
			{row: row, col: col, symbol: '🍊'},
		},
		velRow: 0, velCol: 2,
	})
}

func SplawnBullet2(row, col int) {
	bullets = append(bullets, &GameObject{
		points: []*Point{
			{row: row, col: col, symbol: '🌽'},
		},
		velRow: 0, velCol: 4,
	})
}

func UpdateState() {
	if isGamePaused {
		return
	}

	if score > 200 {
		isGameOver = true
		return
	}
	MoveGameObjekts(append(append(zombies, bullets...), player1, player2))
	MoveGameObjekts2(append(append(zombies2, bullets...), player1, player2))
	MoveGameObjekts3(append(append(zombies3, bullets...), player1, player2))
	UpdateZombies()
	UpdateZombies2()
	UpdateZombies3()
	CollisionDetection()
}

func MoveGameObjekts(objs []*GameObject) {
	for _, obj := range objs {
		for i := range obj.points {
			obj.points[i].row += obj.velRow
			obj.points[i].col += obj.velCol
		}
	}
}

func MoveGameObjekts2(objs []*GameObject) {
	for _, obj := range objs {
		for i := range obj.points {
			obj.points[i].row += obj.velRow
			obj.points[i].col += obj.velCol
		}
	}
}

func MoveGameObjekts3(objs []*GameObject) {
	for _, obj := range objs {
		for i := range obj.points {
			obj.points[i].row += obj.velRow
			obj.points[i].col += obj.velCol
		}
	}
}

func UpdateZombies() {
	spawnChance := rand.Intn(200)
	if spawnChance < 5 {
		SpawnZombie()
	}
}

func SpawnZombie() {
	originRow, originCol := rand.Intn(GameFrameHeight-3), GameFrameWidth-2
	zombies = append(zombies, &GameObject{
		points: []*Point{
			{row: originRow, col: originCol, symbol: '😡'},
			{row: originRow + 1, col: originCol, symbol: '👔'},
			{row: originRow + 2, col: originCol + 1, symbol: '▓'},
			{row: originRow + 3, col: originCol, symbol: '🚒'},
			{row: originRow, col: originCol - 2, symbol: '⛏'},
			{row: originRow, col: originCol + 2, symbol: '🏸'},
			{row: originRow + 1, col: originCol - 1, symbol: '└'},
		},
		velRow: 0, velCol: -2,
	})
}

func UpdateZombies2() {
	spawnChance2 := rand.Intn(200)
	if spawnChance2 < 5 {
		SpawnZombie2()
	}
}

func SpawnZombie2() {
	originRow, originCol := rand.Intn(GameFrameHeight-3), GameFrameWidth-2
	zombies2 = append(zombies2, &GameObject{
		points: []*Point{
			{row: originRow, col: originCol, symbol: '😡'},
			{row: originRow + 1, col: originCol, symbol: '👔'},
			{row: originRow + 2, col: originCol + 1, symbol: '▓'},
			{row: originRow + 3, col: originCol, symbol: '⛸'},
			{row: originRow, col: originCol - 2, symbol: '⛏'},
			{row: originRow + 2, col: originCol - 4, symbol: '🚪'},
			{row: originRow + 1, col: originCol - 4, symbol: '🚪'},
			{row: originRow + 3, col: originCol - 4, symbol: '🚪'},
			{row: originRow, col: originCol - 4, symbol: '🚪'},
			{row: originRow + 1, col: originCol - 1, symbol: '└'},
		},
		velRow: 0, velCol: -1,
		health: 2,
	})
}

func UpdateZombies3() {
	spawnChance3 := rand.Intn(75)
	if spawnChance3 < 5 {
		SpawnZombie3()
	}
}

func SpawnZombie3() {
	originRow, originCol := rand.Intn(GameFrameHeight-3), GameFrameWidth-2
	zombies3 = append(zombies3, &GameObject{
		points: []*Point{
			{row: originRow, col: originCol, symbol: '🎃'},
			{row: originRow + 1, col: originCol, symbol: '👔'},
			{row: originRow + 2, col: originCol + 1, symbol: '▓'},
			{row: originRow + 3, col: originCol, symbol: '⛸'},
			{row: originRow, col: originCol - 2, symbol: '⛏'},
			{row: originRow + 2, col: originCol - 4, symbol: '🚪'},
			{row: originRow + 1, col: originCol - 4, symbol: '🚪'},
			{row: originRow + 3, col: originCol - 4, symbol: '🚪'},
			{row: originRow, col: originCol - 4, symbol: '🚪'},
			{row: originRow + 1, col: originCol - 1, symbol: '└'},
		},
		velRow: 0, velCol: -1,
		health: 3,
	})
}

func CollisionDetection() {
	RemoveObjectsOutOfBounds()
	HandleZombiePlayerCollision()
	HandleZombieBulletCollision()
	RemoveObjectsOutOfBounds2()
	HandleZombiePlayerCollision2()
	HandleZombieBulletCollision2()
	RemoveObjectsOutOfBounds3()
	HandleZombiePlayerCollision3()
	HandleZombieBulletCollision3()
}

///////////////////////////////////////////////////////////

func RemoveObjectsOutOfBounds3() {
	ObjectOutOfBoundsCollision3(zombies3, true, func(idx int) {
		isGameOver = true
	})
	bulletsToRemove := []*GameObject{}
	ObjectOutOfBoundsCollision3(bullets, false, func(idx int) {
		bulletsToRemove = append(bulletsToRemove, bullets[idx])
		bullets = append(bullets[:idx], bullets[idx+1:]...)
	})
	bullets = RemoveGameObjects3(bullets, bulletsToRemove)
}

func HandleZombiePlayerCollision3() {
	for _, z := range zombies3 {
		if AreObjectsColliding3(player1, z, 1) || AreObjectsColliding3(player2, z, 1) {
			isGameOver = true
		}
	}
}

func HandleZombieBulletCollision3() {
	bulletsToRemove := []*GameObject{}
	zombiesToRemove := []*GameObject{}
	for _, b := range bullets {
		for _, z := range zombies3 {
			if AreObjectsColliding3(b, z, 1) {
				bulletsToRemove = append(bulletsToRemove, b)
				z.health--
				if z.health <= 0 {
					zombiesToRemove = append(zombiesToRemove, z)
					score++
				}
				break
			}
		}
	}
	bullets = RemoveGameObjects3(bullets, bulletsToRemove)
	zombies3 = RemoveGameObjects3(zombies3, zombiesToRemove)
}

func AreObjectsColliding3(obj1, obj2 *GameObject, radius int) bool {
	for _, p1 := range obj1.points {
		for _, p2 := range obj2.points {
			if p1.row == p2.row && math.Abs(float64(p1.col-p2.col)) <= float64(radius) {
				return true
			}
		}
	}
	return false
}

func RemoveGameObjects3(source, toRemove []*GameObject) []*GameObject {
	result := []*GameObject{}
	for _, obj1 := range source {
		removed := false
		for _, obj2 := range toRemove {
			if obj1 == obj2 {
				removed = true
				break
			}
		}
		if !removed {
			result = append(result, obj1)
		}
	}
	return result
}

func ObjectOutOfBoundsCollision3(objs []*GameObject, lookAhead bool, callback func(int)) {
	for i, obj := range objs {
		velRow, velCol := obj.velRow, obj.velCol
		if !lookAhead {
			velRow, velCol = 0, 0
		}
		if IsObjectOutOfBounds3(obj, velRow, velCol) {
			callback(i)
		}
	}
}

func IsObjectOutOfBounds3(obj *GameObject, velRow, velCol int) bool {
	for _, p := range obj.points {
		targetRow, targetCol := p.row+velRow, p.col+velCol
		if targetRow < 0 || targetRow >= GameFrameHeight ||
			targetCol < 0 || targetCol >= GameFrameWidth {
			return true
		}
	}
	return false
}

// //////////////////////////////////////////////
func RemoveObjectsOutOfBounds2() {
	ObjectOutOfBoundsCollision2(zombies2, true, func(idx int) {
		isGameOver = true
	})
	bulletsToRemove := []*GameObject{}
	ObjectOutOfBoundsCollision2(bullets, false, func(idx int) {
		bulletsToRemove = append(bulletsToRemove, bullets[idx])
		bullets = append(bullets[:idx], bullets[idx+1:]...)
	})
	bullets = RemoveGameObjects2(bullets, bulletsToRemove)
}

func HandleZombiePlayerCollision2() {
	for _, z := range zombies2 {
		if AreObjectsColliding2(player1, z, 1) || AreObjectsColliding2(player2, z, 1) {
			isGameOver = true
		}
	}
}

func HandleZombieBulletCollision2() {
	bulletsToRemove := []*GameObject{}
	zombiesToRemove := []*GameObject{}
	for _, b := range bullets {
		for _, z := range zombies2 {
			if AreObjectsColliding2(b, z, 1) {
				bulletsToRemove = append(bulletsToRemove, b)
				z.health--
				if z.health <= 0 {
					zombiesToRemove = append(zombiesToRemove, z)
					score++
				}
				break
			}
		}
	}
	bullets = RemoveGameObjects2(bullets, bulletsToRemove)
	zombies2 = RemoveGameObjects2(zombies2, zombiesToRemove)
}

func AreObjectsColliding2(obj1, obj2 *GameObject, radius int) bool {
	for _, p1 := range obj1.points {
		for _, p2 := range obj2.points {
			if p1.row == p2.row && math.Abs(float64(p1.col-p2.col)) <= float64(radius) {
				return true
			}
		}
	}
	return false
}

func RemoveGameObjects2(source, toRemove []*GameObject) []*GameObject {
	result := []*GameObject{}
	for _, obj1 := range source {
		removed := false
		for _, obj2 := range toRemove {
			if obj1 == obj2 {
				removed = true
				break
			}
		}
		if !removed {
			result = append(result, obj1)
		}
	}
	return result
}

func ObjectOutOfBoundsCollision2(objs []*GameObject, lookAhead bool, callback func(int)) {
	for i, obj := range objs {
		velRow, velCol := obj.velRow, obj.velCol
		if !lookAhead {
			velRow, velCol = 0, 0
		}
		if IsObjectOutOfBounds2(obj, velRow, velCol) {
			callback(i)
		}
	}
}

func IsObjectOutOfBounds2(obj *GameObject, velRow, velCol int) bool {
	for _, p := range obj.points {
		targetRow, targetCol := p.row+velRow, p.col+velCol
		if targetRow < 0 || targetRow >= GameFrameHeight ||
			targetCol < 0 || targetCol >= GameFrameWidth {
			return true
		}
	}
	return false
}

// /////////////////////////////////////////////////////////////////
func RemoveObjectsOutOfBounds() {
	ObjectOutOfBoundsCollision(zombies, true, func(idx int) {
		isGameOver = true
	})
	bulletsToRemove := []*GameObject{}
	ObjectOutOfBoundsCollision(bullets, false, func(idx int) {
		bulletsToRemove = append(bulletsToRemove, bullets[idx])
		bullets = append(bullets[:idx], bullets[idx+1:]...)
	})
	bullets = RemoveGameObjects(bullets, bulletsToRemove)
}

func HandleZombiePlayerCollision() {
	for _, z := range zombies {
		if AreObjectsColliding(player1, z, 1) || AreObjectsColliding(player2, z, 1) {
			isGameOver = true
		}
	}
}

func HandleZombieBulletCollision() {
	bulletsToRemove := []*GameObject{}
	zombiesToRemove := []*GameObject{}
	for _, b := range bullets {
		for _, z := range zombies {
			if AreObjectsColliding(b, z, 1) {
				bulletsToRemove = append(bulletsToRemove, b)
				zombiesToRemove = append(zombiesToRemove, z)
				score++
				break
			}
		}
	}
	bullets = RemoveGameObjects(bullets, bulletsToRemove)
	zombies = RemoveGameObjects(zombies, zombiesToRemove)
}

func AreObjectsColliding(obj1, obj2 *GameObject, radius int) bool {
	for _, p1 := range obj1.points {
		for _, p2 := range obj2.points {
			if p1.row == p2.row && math.Abs(float64(p1.col-p2.col)) <= float64(radius) {
				return true
			}
		}
	}
	return false
}

func RemoveGameObjects(source, toRemove []*GameObject) []*GameObject {
	result := []*GameObject{}
	for _, obj1 := range source {
		removed := false
		for _, obj2 := range toRemove {
			if obj1 == obj2 {
				removed = true
				break
			}
		}
		if !removed {
			result = append(result, obj1)
		}
	}
	return result
}

func ObjectOutOfBoundsCollision(objs []*GameObject, lookAhead bool, callback func(int)) {
	for i, obj := range objs {
		velRow, velCol := obj.velRow, obj.velCol
		if !lookAhead {
			velRow, velCol = 0, 0
		}
		if IsObjectOutOfBounds(obj, velRow, velCol) {
			callback(i)
		}
	}
}

func IsObjectOutOfBounds(obj *GameObject, velRow, velCol int) bool {
	for _, p := range obj.points {
		targetRow, targetCol := p.row+velRow, p.col+velCol
		if targetRow < 0 || targetRow >= GameFrameHeight ||
			targetCol < 0 || targetCol >= GameFrameWidth {
			return true
		}
	}
	return false
}

func DrawState() {
	if isGamePaused {
		return
	}
	screen.Clear()
	PrintGameFrame()
	PrintGameFrame2()
	PrintGameObjects(append(append(zombies, bullets...), player1, player2))
	PrintGameObjects2(append(append(zombies2, bullets...), player1, player2))
	PrintGameObjects3(append(append(zombies3, bullets...), player1, player2))
	screen.Show()
}

func PrintGameObjects(objs []*GameObject) {
	for _, obj := range objs {
		for _, p := range obj.points {
			PrintFilledRectInGameFrame(p.row, p.col, 1, 1, p.symbol)
		}
	}
}

func PrintGameObjects2(objs []*GameObject) {
	for _, obj := range objs {
		for _, p := range obj.points {
			PrintFilledRectInGameFrame(p.row, p.col, 1, 1, p.symbol)
		}
	}
}

func PrintGameObjects3(objs []*GameObject) {
	for _, obj := range objs {
		for _, p := range obj.points {
			PrintFilledRectInGameFrame(p.row, p.col, 1, 1, p.symbol)
		}
	}
}

func PrintGameFrame() {
	frameStartRow, frameStartCol := GetGameFrameTopLeft()
	PrintUnfilledRect(
		frameStartRow-1, frameStartCol-1, GameFrameWidth+2, GameFrameHeight+2, GameFrameSymbol)
}

func PrintGameFrame2() {
	frameStartRow, frameStartCol := GetGameFrameTopLeft()
	PrintUnfilledRect(
		frameStartRow-2, frameStartCol-3, GameFrameWidth+4, GameFrameHeight+4, GameFrameSymbol2)
}

func GetGameFrameTopLeft() (int, int) {
	screenWidth, screenHeight := screen.Size()
	return screenHeight/2 - GameFrameHeight/2, screenWidth/2 - GameFrameWidth/2
}

func PrintStringCentered(row, col int, str string) {
	col = col - len(str)/2
	PrintString(row, col, str)
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		PrintFilledRect(row, col, 1, 1, c)
		col += 1
	}
}

func PrintFilledRectInGameFrame(row, col, width, height int, ch rune) {
	frameRow, frameCol := GetGameFrameTopLeft()
	PrintFilledRect(row+frameRow, col+frameCol, width, height, ch)
}

func PrintFilledRect(row, col, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c+1, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func PrintUnfilledRect(row, col, width, height int, ch rune) {
	for c := 0; c < width; c++ {
		screen.SetContent(col+c, row, ch, nil, tcell.StyleDefault)
	}
	for r := 0; r < height-1; r++ {
		screen.SetContent(col, row+r, ch, nil, tcell.StyleDefault)
		screen.SetContent(col+width-1, row+r, ch, nil, tcell.StyleDefault)
	}
	for c := 0; c < width; c++ {
		screen.SetContent(col+c, row+height-1, ch, nil, tcell.StyleDefault)
	}
}
