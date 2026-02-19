[Setup]
AppName=Win Cleaner
AppVersion={#MyAppVersion}
AppPublisher=WinCleaner
AppPublisherURL=https://github.com/your-repo/win-cleaner
DefaultDirName={autopf}\WinCleaner
DefaultGroupName=Win Cleaner
UninstallDisplayIcon={app}\WinCleaner.exe
OutputDir=..\
OutputBaseFilename=WinCleaner-Setup-{#MyAppVersion}
Compression=lzma2/ultra64
SolidCompression=yes
SetupIconFile=..\favicon_io\favicon.ico
ArchitecturesInstallIn64BitMode=x64compatible
WizardStyle=modern
PrivilegesRequired=admin

[Languages]
Name: "chinesesimplified"; MessagesFile: "compiler:Languages\ChineseSimplified.isl"
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "创建桌面快捷方式"; GroupDescription: "附加图标:"

[Files]
Source: "bin\WinCleaner.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\Win Cleaner"; Filename: "{app}\WinCleaner.exe"
Name: "{group}\卸载 Win Cleaner"; Filename: "{uninstallexe}"
Name: "{autodesktop}\Win Cleaner"; Filename: "{app}\WinCleaner.exe"; Tasks: desktopicon

[Run]
Filename: "{app}\WinCleaner.exe"; Description: "启动 Win Cleaner"; Flags: nowait postinstall skipifsilent
