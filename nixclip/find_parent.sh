go_ppid=3264199

# Loop until you find the terminal emulator process
while true; do
    # Get the parent process ID of the current process
    parent_pid=$(ps -o ppid= -p $go_ppid)
    
    # Get the name of the process corresponding to the parent PID
    parent_name=$(ps -o comm= -p $parent_pid)
    
    # Check if the parent process is the terminal emulator process
    if [[ $parent_name == "your_terminal_emulator_name" ]]; then
        # Print the PPID of the terminal emulator process
        echo "PPID of terminal emulator: $parent_pid"
        break
    fi
    
    # Update the current process ID to be the parent process ID
    go_ppid=$parent_pid
done
