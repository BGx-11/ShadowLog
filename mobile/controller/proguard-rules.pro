# ProGuard rules for Controller APK
-keep class com.system.controller.** { *; }
-keepattributes *Annotation*
-keepattributes JavascriptInterface
-keep class * { @android.webkit.JavascriptInterface <methods>; }
