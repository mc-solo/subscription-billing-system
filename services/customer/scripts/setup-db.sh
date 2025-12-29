#!bin/bash

# colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Setting up Customer Service Database...${NC}"


# Check if MySQL is running
if ! mysqladmin ping -h localhost -u root -p 2>/dev/null | grep -q "mysqld is alive"; then
    echo -e "${RED}MySQL is not running. Please start MySQL and try again.${NC}"
    exit 1
fi


# Create database
echo -e "${YELLOW}Creating database...${NC}"
mysql -h localhost -u root -p << EOF
CREATE DATABASE IF NOT EXISTS subscription-billing-sys CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EOF

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Database created successfully!${NC}"
else
    echo -e "${RED}Failed to create database${NC}"
    exit 1
fi

# run migrations
echo -e "${YELLOW}Running migrations...${NC}"
cd migrations
migrate -path . -database "mysql://root:password@tcp(localhost:3306)/subscription-billing-system?multiStatements=true" up
if [$? -eq 0]; then
    echo -e "{GREEN}Migrations completed successfully!${NC}"
else 
    echo -e "${RED}Failed to run migrations${NC}"
    exit 1
fi


# Verify tables were created
echo -e "${YELLOW}Verifying tables...${NC}"
mysql -h localhost -u root -p customer_db << EOF
SHOW TABLES;
SELECT 
    TABLE_NAME,
    TABLE_ROWS,
    DATA_LENGTH,
    INDEX_LENGTH,
    CREATE_TIME
FROM information_schema.TABLES 
WHERE TABLE_SCHEMA = 'customer_db'
ORDER BY TABLE_NAME;
EOF

echo -e "${GREEN}Database setup completed successfully!${NC}"
