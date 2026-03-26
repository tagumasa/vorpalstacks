import 'package:flutter/material.dart';
import 'pages/home/home_page.dart';
import 'pages/services/iam_page.dart';
import 'pages/services/dynamodb_page.dart';
import 'pages/services/sqs_page.dart';
import 'pages/services/sns_page.dart';
import 'pages/services/s3_page.dart';
import 'pages/services/lambda_page.dart';
import 'pages/services/step_functions_page.dart';
import 'pages/services/cognito_page.dart';
import 'pages/services/secrets_manager_page.dart';
import 'pages/services/ssm_page.dart';
import 'pages/services/apigateway_page.dart';
import 'pages/services/cloudfront_page.dart';
import 'pages/services/route53_page.dart';
import 'pages/services/acm_page.dart';
import 'pages/services/cloudwatch_page.dart';
import 'pages/services/athena_page.dart';
import 'pages/services/kms_page.dart';
import 'pages/services/eventbridge_page.dart';
import 'pages/services/scheduler_page.dart';
import 'pages/services/cloudtrail_page.dart';
import 'pages/services/cognitoidentity_page.dart';
import 'pages/services/kinesis_page.dart';
import 'pages/services/ses_page.dart';
import 'pages/services/sts_page.dart';
import 'pages/services/timestream_query_page.dart';
import 'pages/services/timestream_write_page.dart';
import 'pages/services/waf_page.dart';
import 'pages/services/wafv2_page.dart';
import 'pages/settings/settings_page.dart';
import 'widgets/common/sidebar.dart';

class VorpalStackApp extends StatelessWidget {
  const VorpalStackApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'VorpalStack',
      debugShowCheckedModeBanner: false,
      theme: _buildLightTheme(),
      darkTheme: _buildDarkTheme(),
      themeMode: ThemeMode.light,
      home: const MainLayout(),
    );
  }

  ThemeData _buildLightTheme() {
    const primaryBlue = Color(0xFF232F3E);
    const accentBlue = Color(0xFF2196F3);
    const lightGray = Color(0xFFF2F3F4);
    const mediumGray = Color(0xFFE5E5E5);
    const darkGray = Color(0xFF545B64);
    const textPrimary = Color(0xFF16191F);

    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.light(
        primary: primaryBlue,
        secondary: accentBlue,
        surface: Colors.white,
        error: const Color(0xFFD32F2F),
        onPrimary: Colors.white,
        onSecondary: Colors.white,
        onSurface: textPrimary,
        outline: mediumGray,
      ),
      scaffoldBackgroundColor: lightGray,
      appBarTheme: const AppBarTheme(
        backgroundColor: primaryBlue,
        foregroundColor: Colors.white,
        elevation: 0,
        titleTextStyle: TextStyle(
          color: Colors.white,
          fontSize: 18,
          fontWeight: FontWeight.w500,
        ),
      ),
      cardTheme: CardThemeData(
        color: Colors.white,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
          side: const BorderSide(color: mediumGray),
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: accentBlue,
          foregroundColor: Colors.white,
          elevation: 0,
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(4)),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(4),
          borderSide: const BorderSide(color: mediumGray),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(4),
          borderSide: const BorderSide(color: primaryBlue, width: 2),
        ),
        filled: true,
        fillColor: Colors.white,
        contentPadding: const EdgeInsets.symmetric(
          horizontal: 12,
          vertical: 12,
        ),
      ),
      dividerColor: mediumGray,
      iconTheme: const IconThemeData(color: darkGray),
    );
  }

  ThemeData _buildDarkTheme() {
    const darkBg = Color(0xFF0F1B2A);
    const darkSurface = Color(0xFF16202E);
    const darkAccent = Color(0xFF2196F3);
    const darkOrange = Color(0xFFFF9900);
    const textLight = Color(0xFFE5E5E5);
    const borderDark = Color(0xFF2D3F52);

    return ThemeData(
      useMaterial3: true,
      colorScheme: const ColorScheme.dark(
        primary: darkAccent,
        secondary: darkOrange,
        surface: darkSurface,
        error: Color(0xFFEF5350),
        onPrimary: Colors.white,
        onSecondary: Colors.white,
        onSurface: textLight,
        outline: borderDark,
      ),
      scaffoldBackgroundColor: darkBg,
      cardTheme: CardThemeData(
        color: darkSurface,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
          side: const BorderSide(color: borderDark),
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: darkAccent,
          foregroundColor: Colors.white,
          elevation: 0,
          padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(4)),
        ),
      ),
    );
  }
}

class MainLayout extends StatefulWidget {
  const MainLayout({super.key});

  @override
  State<MainLayout> createState() => _MainLayoutState();
}

class _MainLayoutState extends State<MainLayout> {
  int _selectedIndex = 0;

  final List<Widget> _pages = [
    const HomePage(),
    const IAMPage(),
    const DynamodbPage(),
    const SqsPage(),
    const SnsPage(),
    const S3Page(),
    const LambdaPage(),
    const StepFunctionsPage(),
    const CognitoPage(),
    const SecretsManagerPage(),
    const SsmPage(),
    const ApigatewayPage(),
    const CloudfrontPage(),
    const Route53Page(),
    const AcmPage(),
    const CloudwatchPage(),
    const AthenaPage(),
    const KMSPage(),
    const EventBridgePage(),
    const SchedulerPage(),
    const CloudTrailPage(),
    const CognitoIdentityPage(),
    const KinesisPage(),
    const SESPage(),
    const STSPage(),
    const TimestreamQueryPage(),
    const TimestreamWritePage(),
    const WAFPage(),
    const WAFv2Page(),
    const SettingsPage(),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Column(
        children: [
          _buildTopBar(context),
          Expanded(
            child: Row(
              children: [
                Sidebar(
                  selectedIndex: _selectedIndex,
                  onTap: (index) => setState(() => _selectedIndex = index),
                ),
                Expanded(child: _pages[_selectedIndex]),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTopBar(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return Container(
      height: 56,
      decoration: BoxDecoration(
        color: cs.surface,
        border: Border(bottom: BorderSide(color: cs.outline)),
      ),
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Row(
        children: [
          Expanded(
            child: Container(
              constraints: const BoxConstraints(maxWidth: 400),
              child: TextField(
                decoration: InputDecoration(
                  hintText: 'Search services, resources...',
                  prefixIcon: Icon(
                    Icons.search,
                    size: 20,
                    color: cs.onSurface.withValues(alpha: 0.5),
                  ),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(6),
                    borderSide: BorderSide(color: cs.outline),
                  ),
                  contentPadding: const EdgeInsets.symmetric(
                    horizontal: 12,
                    vertical: 8,
                  ),
                  filled: true,
                  fillColor: cs.surface,
                ),
              ),
            ),
          ),
          const Spacer(),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
            decoration: BoxDecoration(
              border: Border.all(color: cs.outline),
              borderRadius: BorderRadius.circular(4),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(
                  Icons.public,
                  size: 16,
                  color: cs.onSurface.withValues(alpha: 0.7),
                ),
                const SizedBox(width: 8),
                Text(
                  'us-east-1',
                  style: TextStyle(
                    fontSize: 13,
                    color: cs.onSurface,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                Icon(
                  Icons.arrow_drop_down,
                  size: 20,
                  color: cs.onSurface.withValues(alpha: 0.7),
                ),
              ],
            ),
          ),
          const SizedBox(width: 12),
          IconButton(
            icon: Badge(
              label: const Text('3'),
              child: Icon(
                Icons.notifications_outlined,
                color: cs.onSurface.withValues(alpha: 0.7),
              ),
            ),
            onPressed: () {},
          ),
          const SizedBox(width: 8),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: cs.primary.withValues(alpha: 0.1),
              borderRadius: BorderRadius.circular(4),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                CircleAvatar(
                  radius: 14,
                  backgroundColor: cs.primary,
                  child: Text(
                    'U',
                    style: TextStyle(
                      color: cs.onPrimary,
                      fontSize: 12,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                Text(
                  'test-user',
                  style: TextStyle(
                    fontSize: 13,
                    color: cs.onSurface,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                Icon(
                  Icons.arrow_drop_down,
                  size: 20,
                  color: cs.onSurface.withValues(alpha: 0.7),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
