# Examples 



# Upadating user

```go
update := chawk.UserUpdate{
    Name:    &chawk.NameUpdate{ Family: pointer.ToString("Smith") },
    Contact: &chawk.ContactUpdate{ Email: pointer.ToString("jsmith@univ.edu") },
    InstitutionRoleIDs: []string{"FACULTY"},
}

err := client.Users.Update(ctx, "jdoe", update)


or 
err = client.Users.UpdateEmail(ctx, "jdoe", "jsmith@univ.edu")

```