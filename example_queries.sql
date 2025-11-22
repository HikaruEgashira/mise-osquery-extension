-- Query all packages
SELECT * FROM version_manager_packages;

-- Query packages by tool
SELECT tool, version, manager FROM version_manager_packages WHERE tool = 'node';

-- Count packages per manager
SELECT manager, COUNT(*) as count FROM version_manager_packages GROUP BY manager;

-- Find specific tool versions
SELECT * FROM version_manager_packages WHERE tool LIKE '%node%';

-- List unique tools
SELECT DISTINCT tool FROM version_manager_packages ORDER BY tool;

-- Find all Python installations
SELECT tool, version, install_path FROM version_manager_packages WHERE tool = 'python';

-- Find all Go installations
SELECT tool, version, install_path FROM version_manager_packages WHERE tool = 'go';

-- Count total number of installed tools
SELECT COUNT(*) as total_tools FROM version_manager_packages;

-- Count tools per version manager
SELECT manager, tool, COUNT(*) as version_count
FROM version_manager_packages
GROUP BY manager, tool
ORDER BY manager, tool;

-- Find tools with multiple versions installed
SELECT tool, COUNT(*) as version_count, GROUP_CONCAT(version, ', ') as versions
FROM version_manager_packages
GROUP BY tool
HAVING COUNT(*) > 1
ORDER BY version_count DESC;

-- Find specific version of a tool
SELECT * FROM version_manager_packages WHERE tool = 'node' AND version LIKE '20%';

-- List all tools installed via mise
SELECT * FROM version_manager_packages WHERE manager = 'mise' ORDER BY tool;

-- List all tools installed via asdf
SELECT * FROM version_manager_packages WHERE manager = 'asdf' ORDER BY tool;

-- Find installation paths for a specific tool
SELECT tool, version, install_path FROM version_manager_packages WHERE tool = 'ruby';

-- Check if a specific tool version is installed
SELECT EXISTS(
  SELECT 1 FROM version_manager_packages
  WHERE tool = 'node' AND version = '20.10.0'
) as is_installed;
