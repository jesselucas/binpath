## binpath
Finds the nearest /bin and executes a command in it.

## Usage
`binpath command [arguments]`

## Bash Complete
If you'd like to enable bash complete add `bp_bash_complete` to your file system (such as your home directory) and in your .bashrc add it as a source:

Ex. `source $HOME/bp_bash_complete`

### Alias
`binpath` bash complete is setup to use the alias `bp`. 

Add the alias to your .bashrc: 
Ex. `alias bp='binpath'`