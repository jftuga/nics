
$ErrorActionPreference = 'Stop';
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64      = 'https://github.com/jftuga/nics/releases/download/v1.4.2/nics_1.4.2_windows_amd64.zip'

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  unzipLocation = $toolsDir
  fileType      = 'exe'
  url           = $url
  url64bit      = $url64

  softwareName  = 'nics*'
  checksum64    = '1d9726a70f1be3f6a81606ae010e730ae5d66495fb50cb07567bed8a77488ff8'
  checksumType64= 'sha256'

  validExitCodes= @(0)
}

Install-ChocolateyZipPackage @packageArgs
