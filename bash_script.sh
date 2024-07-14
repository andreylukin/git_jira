function git-jira() {
    # Check if the argument is provided
    if [ -z "$1" ]; then
        echo "Error: No string argument provided."
        return 1
    fi

    # Call the Go script to get the branch name
    # Capture all output
    go_output=$(go run /home/alukin/repos/git_jira/main.go "$1")
    if [ $? -ne 0 ]; then
        echo "Error: Failed to generate branch name using the Go script."
        return 1
    fi

    # Split the output into two lines
    IFS=$'\n' read -r user_message branch_name <<< "$go_output"

    # Print the first line for the user
    echo "$user_message"    

    # Validate branch name
    if [ -z "$branch_name" ]; then
        echo "Error: Go script did not output a valid branch name."
        return 1
    fi

    # Get the default branch name
    default_branch=$(git remote show origin | sed -n '/HEAD branch/s/.*: //p')
    if [ $? -ne 0 ]; then
        echo "Error: Unable to determine the default branch."
        return 1
    fi

    # Fetch the latest changes from the remote
    git fetch origin
    if [ $? -ne 0 ]; then
        echo "Error: Failed to fetch the latest changes from the remote."
        return 1
    fi

    # Create a new branch off the default branch and switch to it
    git checkout -b "$branch_name" "origin/$default_branch"
    if [ $? -eq 0 ]; then
        echo "Created and switched to new branch '$branch_name' based on '$default_branch'"
    else
        echo "Error: Failed to create and switch to new branch '$branch_name'"
        return 1
    fi
}