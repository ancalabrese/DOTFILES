function adbReset --wraps='adb shell am broadcast -a android.intent.action.MASTER_CLEAR' --description 'alias adbReset=adb shell am broadcast -a android.intent.action.MASTER_CLEAR'
  adb shell am broadcast -a android.intent.action.MASTER_CLEAR $argv
        
end
