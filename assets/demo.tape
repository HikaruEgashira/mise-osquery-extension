Output demo.gif

Set FontSize 24
Set Width 1200
Set Height 600
Set Padding 20
Set TypingSpeed 30ms

Type "osqueryi --extension $(which version_manager_packages_extension)" Sleep 200ms Enter

Sleep 1.5s

Type "SELECT tool, version, manager FROM version_manager_packages LIMIT 5;" Sleep 200ms Enter

Sleep 2s

Type "SELECT manager, COUNT(*) as count FROM version_manager_packages GROUP BY manager;" Sleep 200ms Enter

Sleep 2s

Type "SELECT tool, version FROM version_manager_packages WHERE tool = 'node';" Sleep 200ms Enter

Sleep 2s

Type ".quit" Sleep 200ms Enter

Sleep 500ms
