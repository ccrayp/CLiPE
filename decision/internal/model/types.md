| Type          | Описание                      | Операторы                     | Источник (система / PAM)        | Как получить                         |
|---------------|------------------------------|--------------------------------|----------------------------------|--------------------------------------|
| gid           | Основная группа              | equals, not_equals             | libc                             | getgid()                             |
| group         | Группа пользователя          | equals, contains               | libc                             | getgroups()                          |
| ip            | IP клиента                   | equals, in, not_in             | PAM / socket                     | pam_get_item(PAM_RHOST)              |
| hostname      | Имя хоста                    | equals, not_equals             | libc                             | gethostname()                        |
| time          | Текущее время                | between                        | libc                             | time(), localtime()                  |
| weekday       | День недели                  | in, equals                     | libc                             | localtime()                          |