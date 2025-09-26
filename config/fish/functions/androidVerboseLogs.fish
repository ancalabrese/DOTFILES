function androidVerboseLogs --wraps='adb shell setprop log.tag.all VERBOSE' --description 'alias androidVerboseLogs adb shell setprop log.tag.all VERBOSE'
  adb shell setprop log.tag.all VERBOSE $argv
        
end
