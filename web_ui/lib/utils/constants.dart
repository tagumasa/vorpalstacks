import 'package:flutter/material.dart';

class AppColors {
  AppColors._();

  static const Color primaryBlue = Color(0xFF232F3E);
  static const Color accentBlue = Color(0xFF2196F3);
  static const Color lightBlue = Color(0xFF64B5F6);
  static const Color infoBlue = Color(0xFF0288D1);
  static const Color awsOrange = Color(0xFFFF9900);

  static const Color darkGray = Color(0xFF545B64);
  static const Color mediumGray = Color(0xFF879596);
  static const Color lightGray = Color(0xFFE5E5E5);
  static const Color backgroundGray = Color(0xFFF2F3F4);

  static const Color success = Color(0xFF4CAF50);
  static const Color warning = Color(0xFFFF9800);
  static const Color error = Color(0xFFD32F2F);
  static const Color info = Color(0xFF2196F3);

  static const Color darkBackground = Color(0xFF0F1B2A);
  static const Color darkSurface = Color(0xFF16202E);
  static const Color darkBorder = Color(0xFF2D3F52);
}

class AppSpacing {
  AppSpacing._();

  static const double xs = 4.0;
  static const double sm = 8.0;
  static const double md = 16.0;
  static const double lg = 24.0;
  static const double xl = 32.0;
  static const double xxl = 48.0;
}

class AppTextStyles {
  AppTextStyles._();

  static const TextStyle heading1 = TextStyle(
    fontSize: 28,
    fontWeight: FontWeight.w600,
  );

  static const TextStyle heading2 = TextStyle(
    fontSize: 24,
    fontWeight: FontWeight.w600,
  );

  static const TextStyle heading3 = TextStyle(
    fontSize: 20,
    fontWeight: FontWeight.w600,
  );

  static const TextStyle body = TextStyle(
    fontSize: 14,
    fontWeight: FontWeight.w400,
  );

  static const TextStyle bodyBold = TextStyle(
    fontSize: 14,
    fontWeight: FontWeight.w600,
  );

  static const TextStyle caption = TextStyle(
    fontSize: 12,
    fontWeight: FontWeight.w400,
  );

  static const TextStyle label = TextStyle(
    fontSize: 11,
    fontWeight: FontWeight.w500,
    letterSpacing: 0.5,
  );
}

class AppBreakpoints {
  AppBreakpoints._();

  static const double desktop = 1200.0;
  static const double tablet = 768.0;
  static const double mobile = 0.0;

  static const double sidebarWidthDesktop = 250.0;
  static const double sidebarWidthTablet = 200.0;
}
