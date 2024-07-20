# What does it do?
Basically it helps you find what DB tables are used in a particular Spring Boot project.
It supposes that `Hibernate` is used and entities have such annotation: `@Table(name = "table_name")`

# Examples of the commands
1. `table-sleuth s2t -p /home/spring-boot-project1 -p /home/spring-boot-project2`
2. `table-sleuth t2s -d /home/spring-boot-projects-dir`
