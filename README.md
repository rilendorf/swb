# swb - simple web blog

This is a simple web blog

It generates static .html files from .md files.

---

## Usage

1. install golang

2. install swb `go install github.com/DerZombiiie/swb`

3. install pandoc so `apt install pandoc` or `pacman -Sy pandoc` or whatever is applicable for you

4. create blog `swb create <path>`
   path will be created

5. make `swb make <path>`
   !may though pandoc warnings which is OK!

6. check output at `<path>/out/`
