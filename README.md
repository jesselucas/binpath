## binpath
Looks for the nearest `/bin` folder and executes a command in it. First searches your current folder and then searches up until it finds a `/bin` folder.

## Usage
`binpath command [arguments]`

## Options
`-list, -ls`
  list directory contents of nearest `/bin` folder

## Bash Complete
If you'd like to enable bash complete add the `bp_bash_complete` file to your file system (such as your home directory). Edit your `.bashrc` to add it as a source:

Ex. `source $HOME/bp_bash_complete`

### Alias
`binpath` bash complete is setup to use the alias `bp`. 

Add the alias to your .bashrc: 
Ex. `alias bp='binpath'`