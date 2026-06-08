package main

import (
	"log"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
)

type model struct {
	cursor int
	page   string
	frame  int
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

const (
	pageHome     = "home"
	pageProjects = "projects"
	pageContact  = "contact"
	pageSocials  = "socials"
)

var art string
var logo string

var logoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#809abf")). // #809abf
	Bold(true)

var grey = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#8f8f8f")). // #8f8f8f
	Bold(true)

const hostKey = "id_ed25519"

func loadAssets() {
	var err error

	artBytes, err := os.ReadFile("art.txt")
	if err != nil {
		log.Fatal("failed to load art.txt:", err)
	}
	art = string(artBytes)

	logoBytes, err := os.ReadFile("logo.txt")
	if err != nil {
		log.Fatal("failed to load logo.txt:", err)
	}
	logo = string(logoBytes)
}

func initialModel() model {
	return model{
		page: pageHome,
	}
}

/*
	func ensureHostKey() {
		if _, err := os.Stat(hostKey); err == nil {
			return
		}

		cmd := exec.Command(
			"ssh-keygen",
			"-t", "ed25519",
			"-f", hostKey,
			"-N", "",
		)

		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("ssh-keygen failed: %v\n%s", err, out)
		}
	}
*/
func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tickMsg:
		m.frame++
		return m, tick()

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "left":
			if m.page == pageHome && m.cursor > 0 {
				m.cursor--
			}

		case "right":
			if m.page == pageHome && m.cursor < 2 {
				m.cursor++
			}

		case "enter":
			if m.page == pageHome {
				switch m.cursor {
				case 0:
					m.page = pageProjects
				case 1:
					m.page = pageContact
				case 2:
					m.page = pageSocials
				}
			}

		case "esc":
			m.page = pageHome
			m.cursor = 0
		}
	}

	return m, nil
}

func (m model) View() string {

	switch m.page {

	// ---------------- HOME ----------------
	case pageHome:

		// --- sparkles ---
		sparkles := []rune{'✦', '✧', '⋆', '✷'}

		lines := strings.Split(logo, "\n")

		for i := range lines {
			if len(lines[i]) == 0 {
				continue
			}

			runes := []rune(lines[i])

			if len(runes) > 5 {
				runes[(m.frame*7+i*3)%len(runes)] = sparkles[(m.frame+i)%len(sparkles)]
			}
			if len(runes) > 20 {
				runes[(m.frame*11+i*5)%len(runes)] = sparkles[(m.frame+i+1)%len(sparkles)]
			}
			if len(runes) > 40 {
				runes[(m.frame*13+i*7)%len(runes)] = sparkles[(m.frame+i+2)%len(sparkles)]
			}

			lines[i] = string(runes)
		}

		animatedLogo := logoStyle.Render(strings.Join(lines, "\n"))

		// --- RIGHT SIDE TEXT ---
		var right strings.Builder

		right.WriteString(animatedLogo + "\n\n")

		right.WriteString("(Tahmid) is a prospective computer science and political science student from the NY metro area.\n")
		right.WriteString("He is currently a high school student in 11th grade and is working on projects in his freetime\n")
		right.WriteString("He develops for several platforms such as the 3DS, desktop, mobile, and many more.\n\n")

		right.WriteString(grey.Render("lately, tahmid is passionate about maintaining creativity especially in a not so creative industry") + "\n")
		right.WriteString(grey.Render("you can often find him thinking of what to make (to see what he's made so far check projects!)") + "\n")
		right.WriteString(grey.Render("his work often attempts to be a bit more fun and creative in a world that might not be cut out for him") + "\n\n")

		right.WriteString(grey.Render("\nlearn more using the menus below!"))
		right.WriteString(grey.Render("\nthanks for stopping by! we hope to see you again soon! \n"))

		// --- MENU ---
		menu := []string{"Projects", "Contact", "Socials"}

		right.WriteString("\n\n")
		for i := range menu {
			if i == m.cursor {
				menu[i] = logoStyle.Render("• " + menu[i])
			}
		}

		right.WriteString(strings.Join(menu, "   "))

		// --- LEFT SIDE ART ---
		left := art + "\n"
		right.WriteString(grey.Render("\n\n[← → to navigate · enter to select · q to quit]"))

		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			left,
			"    ",
			right.String(),
		)

	// projects
	case pageProjects:
		var s strings.Builder

		s.WriteString(logoStyle.Render("\nprojects"))
		s.WriteString("\n──────────────────\n\n")

		s.WriteString(logoStyle.Render("mp3ds"))
		s.WriteString(" - a mp3 player for your 3ds\n\n")

		s.WriteString(logoStyle.Render("ssh portfolio"))
		s.WriteString(" - the portfolio you're seeing now!\n\n")

		s.WriteString(logoStyle.Render("bangaranga dimension"))
		s.WriteString(" - a new dimension for minecraft\n\n")

		s.WriteString("for a full list of my projects (+ more details) see my github @")
		s.WriteString(logoStyle.Render(" www.github.com/tr1vee"))

		s.WriteString(grey.Render("\n\n[esc to go back]"))

		return s.String()

	// contact page
	case pageContact:

		var s strings.Builder

		s.WriteString(logoStyle.Render("\ncontact"))
		s.WriteString("\n──────────────────\n\n")

		s.WriteString(logoStyle.Render("email"))
		s.WriteString(": tahmid.raahil0210@gmail.com\n\n")

		s.WriteString(logoStyle.Render("linkedin"))
		s.WriteString(": www.linkedin.com/in/tahmid-raahil-tr1\n\n")

		s.WriteString(logoStyle.Render("discord"))
		s.WriteString(": tr1ve\n\n")

		s.WriteString(grey.Render("[esc to go back]"))

		return s.String()

	// social page
	case pageSocials:
		return logoStyle.Render("\nsocials") + `
──────────────────
` + logoStyle.Render("tiktok") + `: tr11vee` +

			logoStyle.Render("\n\nroblox") + `: tr1ve/bohosalno8` +

			logoStyle.Render("\n\nx/twitter") + `: @tr11vee

` + grey.Render("[esc to go back]")

	default:
		return "loading..."
	}

}

/*func (m model) View() string {
	return "SSH WORKING"
} */

func main() {
	log.Println("starting")

	loadAssets()
	log.Printf("art length=%d", len(art))
	log.Printf("logo length=%d", len(logo))
	// ensureHostKey()

	s, err := wish.NewServer(
		wish.WithAddress(":23234"),
		wish.WithHostKeyPath(hostKey),
		wish.WithMiddleware(
			bm.Middleware(func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
				pty, _, active := sess.Pty()

				if !active {
					log.Println("env:", sess.Environ())

				}
				log.Println("PTY active:", active)
				log.Println("TERM:", pty.Term)
				log.Printf("TERM=%s", sess.Environ())
				return initialModel(), []tea.ProgramOption{
					tea.WithAltScreen(),
				}
			}),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Running SSH server on :23234")
	log.Println("TERM =", os.Getenv("TERM"))

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
