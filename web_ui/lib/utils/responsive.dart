import 'package:flutter/material.dart';

class Responsive extends StatelessWidget {
  final Widget mobile;
  final Widget? tablet;
  final Widget desktop;

  const Responsive({
    super.key,
    required this.mobile,
    this.tablet,
    required this.desktop,
  });

  static bool isMobile(BuildContext context) =>
      MediaQuery.of(context).size.width < 768;

  static bool isTablet(BuildContext context) {
    final w = MediaQuery.of(context).size.width;
    return w >= 768 && w < 1200;
  }

  static bool isDesktop(BuildContext context) =>
      MediaQuery.of(context).size.width >= 1200;

  @override
  Widget build(BuildContext context) {
    final w = MediaQuery.of(context).size.width;
    if (w >= 1200) return desktop;
    if (w >= 768) return tablet ?? desktop;
    return mobile;
  }
}

class ResponsiveLayout {
  static int getGridCrossAxisCount(BuildContext context) {
    if (Responsive.isDesktop(context)) return 4;
    if (Responsive.isTablet(context)) return 2;
    return 1;
  }

  static EdgeInsets getPadding(BuildContext context) {
    if (Responsive.isDesktop(context)) return const EdgeInsets.all(24);
    if (Responsive.isTablet(context)) return const EdgeInsets.all(16);
    return const EdgeInsets.all(12);
  }

  static double getSidebarWidth(BuildContext context) {
    if (Responsive.isDesktop(context)) return 250;
    if (Responsive.isTablet(context)) return 200;
    return 0;
  }
}
