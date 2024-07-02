Create a database (First time)
-----------------------------------------------------

#### Install MySQL on Mac

```bash
# Install
brew install mysql

# Start
brew services start mysql
```

#### Access MySQL
```bash
mysql -u root
```

Add User & DB for user

```sql
CREATE DATABASE gfly CHARACTER SET utf8mb4;
CREATE USER 'user'@'%' IDENTIFIED BY 'secret';
GRANT ALL PRIVILEGES ON gfly.* TO 'user'@'%';
FLUSH PRIVILEGES;
```
