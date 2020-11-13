package nc

import (
	"log"
	"os"
	"path"
)

var timeFormat = "02-01-2006 15:04:05"

const ( // response
	startupMessage   = "Starting on port %v\n"
	yourName         = "[PLEASE ENTER NAME]: "
	groupName        = "[PLEASE ENTER GROUP NAME]: "
	joinGroup        = "[WHICH GROUP TO JOIN]: "
	leftRoom         = "[%v HAS LEFT THE GROUP]\n"
	notInThisGroup   = "[CURRENTLY, YOUR NOT IN THIS GROUP]\n"
	notInGroup       = "[YOU ARE NOT A MEMBER OF THIS GROUP]\n"
	noSuchGroup      = "[NO SUCH GROUP: %v]\n"
	overGroupLimit   = "[CANNOT JOIN %v. LIMIT HAS REACHED]\n"
	joinedGroup      = "[%v HAS JOINED GROUP]\n"
	inThisGroup      = "[ALREADY MEMBER OF THIS GROUP]\n"
	leaveFirst       = "[MEMBER OF ANOTHER GROUP. LEAVING %v]\n"
	leavingGroup     = "[LEAVING GROUP %v]\n"
	invalidGroupName = "[INVALID GROUP NAME %v.]\n"
	internalError    = "[INTERNAL SERVER ERROR.\n" // ERROR 505 INTERNAL SERVER ERROR]
	groupExists      = "[GROUP %v ALREADY EXISTS.]\n"
	existingUsername = "[EXISTING USERNAME %v. PLEASE, TRY ANOTHER ONE.]\n"
	changedName      = "[User '%v' HAS CHANGED NAME TO '%v']\n"
)

const ( // commands
	quit        = "--quit"
	groups      = "--groups"
	join        = "--join"
	create      = "--create"
	leave       = "--leave"
	deleteGroup = "--deleteGroup"
	changeName  = "--rename"
)

const ( // config
	groupLimit = 20
)

const linuxIcon = "" +
	"Welcome to TCP-Chat!\n" +
	"         _nnnn_\n" +
	"        dGGGGMMb\n" +
	"       @p~qp~~qMb\n" +
	"       M|@||@) M|\n" +
	"       @,----.JM|\n" +
	"      JS^\\__/  qKL\n" +
	"     dZP        qKRb\n" +
	"    dZP          qKKb\n" +
	"   fZP            SMMb\n" +
	"   HZM            MMMM\n" +
	"   FqM            MMMM\n" +
	" __| \".        |\\dS\"qML\n" +
	" |    `.       | `' \\Zq\n" +
	"_)      \\.___.,|     .'\n" +
	"\\____   )MMMMMP|   .'\n" +
	"     `-'       `--'\n"

func InitLogger() (*os.File, error) {
	file, err := os.OpenFile(path.Join("..", "logs", "info.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Print("###################################")
	return file, err
}
