package gui

import
import "github.com/therecipe/qt/widgets"

{
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/uitools"
	"github.com/therecipe/qt/widgets"
}

var actionZMX_Settings *widgets.QAction
var	actionOffline_Wiki *widgets.QAction
var actionZimbra_Support *widgets.QAction
var actionOfficial_Documentation *widgets.QAction
var actionCopy *widgets.QAction
var actionCut widgets.QAction
var actionPaste widgets.QAction
var actionOpen_Settings_Backup *widgets.QAction
var actionSave_Settings_Backup *widgets.QAction
var actionSync_Settings widgets.QAction
var actionOpen_zmX_Log *widgets.QAction
var actionAbout_zmX *widgets.QAction
var actionOpen_Remote_Support_Session *widgets.QAction
var	centralwidget *widgets.QWidget
var	gridLayout_4 *widgets.QGridLayout
var	mainTabs *widgets.QTabWidget
var	tab *widgets.QWidget
var	gridLayout *widgets.QGridLayout
var	groupBox_3 *widgets.QGroupBox
var	textBrowser_3 *widgets.QTextBrowser
var	groupBox *widgets.QGroupBox
var	textBrowser *widgets.QTextBrowser
var	groupBox_2 *widgets.QGroupBox
var	textBrowser_2 *widgets.QTextBrowser
var	MetadataScanTab *widgets.QWidget
var	gridLayout_8 *widgets.QGridLayout
var	label *widgets.QLabel
var	checkBox *widgets.QCheckBox
var	comboBox_2 *widgets.QComboBox
var	pushButton *widgets.QPushButton
var	comboBox *widgets.QComboBox
var	tableWidget *widgets.QTableWidget
var	menubar *widgets.QMenuBar
var	menuFile *widgets.QMenu
var	menuEdit *widgets.QMenu
var	menuLibrary *widgets.QMenu
var	menuSettings *widgets.QMenu
var	statusbar *widgets.QStatusBar
var	dockWidget_2 *widgets.QDockWidget
var	dockWidgetContents_2 *widgets.QWidget
var	gridLayout_2 *widgets.QGridLayout
var	toolBox *widgets.QToolBox
var	page *widgets.QWidget
var	gridLayout_3 *widgets.QGridLayout
var	treeWidget *widgets.QTreeWidget
var	page_3 *widgets.QWidget
var	page_4 *widgets.QWidget
var	page_2 *widgets.QWidget
var	gridLayout_5 *widgets.QGridLayout
var	reportsDiagTree *widgets.QTreeWidget
var	dockWidget_3 *widgets.QDockWidget
var	dockWidgetContents_4 *widgets.QWidget
var	verticalLayout_2 *widgets.QVBoxLayout
var	tabWidget_2 *widgets.QTabWidget
var	tab_3 *widgets.QWidget
var	tab_3 *widgets.QWidget
var	verticalLayout *widgets.QVBoxLayout
var	zmxlogger *widgets.QTextEdit
var	tab_4 *widgets.QWidget
var	tab_5 *widgets.QWidget


//	QWidget *centralwidget;
//	QGridLayout *gridLayout_4;
//	QTabWidget *mainTabs;
//	QWidget *tab;
//	QGridLayout *gridLayout;
//	QGroupBox *groupBox_3;
//	QTextBrowser *textBrowser_3;
//	QGroupBox *groupBox;
//	QTextBrowser *textBrowser;
//	QGroupBox *groupBox_2;
//	QTextBrowser *textBrowser_2;
//	QWidget *MetadataScanTab;
//	QGridLayout *gridLayout_8;
//	QLabel *label;
//	QCheckBox *checkBox;
//	QComboBox *comboBox_2;
//	QPushButton *pushButton;
//	QComboBox *comboBox;
//	QTableWidget *tableWidget;
//	QMenuBar *menubar;
//	QMenu *menuFile;
//	QMenu *menuEdit;
//	QMenu *menuLibrary;
//	QMenu *menuSettings;
//	QStatusBar *statusbar;
//	QDockWidget *dockWidget_2;
//	QWidget *dockWidgetContents_2;
//	QGridLayout *gridLayout_2;
//	QToolBox *toolBox;
//	QWidget *page;
//	QGridLayout *gridLayout_3;
//	QTreeWidget *treeWidget;
//	QWidget *page_3;
//	QWidget *page_4;
//	QWidget *page_2;
//	QGridLayout *gridLayout_5;
//	QTreeWidget *reportsDiagTree;
//	QDockWidget *dockWidget_3;
//	QWidget *dockWidgetContents_4;
//	QVBoxLayout *verticalLayout_2;
//	QTabWidget *tabWidget_2;
//	QWidget *tab_3;
//	QVBoxLayout *verticalLayout;
//	QTextEdit *zmxlogger;
//	QWidget *tab_4;
//	QWidget *tab_5;
func CreateMainWindow(){
	//
	//class Ui_MainWindow
	//{
	//public:

	//
	//	void setupUi(QMainWindow *MainWindow)
	//	{
	//		if (MainWindow->objectName().isEmpty())
	//			MainWindow->setObjectName(QString::fromUtf8("MainWindow"));
	//		MainWindow->resize(1151, 742);
	//		actionZMX_Settings = new QAction(MainWindow);
	//		actionZMX_Settings->setObjectName(QString::fromUtf8("actionZMX_Settings"));
	//		actionOffline_Wiki = new QAction(MainWindow);
	//		actionOffline_Wiki->setObjectName(QString::fromUtf8("actionOffline_Wiki"));
	//		actionZimbra_Support = new QAction(MainWindow);
	//		actionZimbra_Support->setObjectName(QString::fromUtf8("actionZimbra_Support"));
	//		actionOfficial_Documentation = new QAction(MainWindow);
	//		actionOfficial_Documentation->setObjectName(QString::fromUtf8("actionOfficial_Documentation"));
	//		actionCopy = new QAction(MainWindow);
	//		actionCopy->setObjectName(QString::fromUtf8("actionCopy"));
	//		actionCut = new QAction(MainWindow);
	//		actionCut->setObjectName(QString::fromUtf8("actionCut"));
	//		actionPaste = new QAction(MainWindow);
	//		actionPaste->setObjectName(QString::fromUtf8("actionPaste"));
	//		actionOpen_Settings_Backup = new QAction(MainWindow);
	//		actionOpen_Settings_Backup->setObjectName(QString::fromUtf8("actionOpen_Settings_Backup"));
	//		actionSave_Settings_Backup = new QAction(MainWindow);
	//		actionSave_Settings_Backup->setObjectName(QString::fromUtf8("actionSave_Settings_Backup"));
	//		actionSync_Settings = new QAction(MainWindow);
	//		actionSync_Settings->setObjectName(QString::fromUtf8("actionSync_Settings"));
	//		actionOpen_zmX_Log = new QAction(MainWindow);
	//		actionOpen_zmX_Log->setObjectName(QString::fromUtf8("actionOpen_zmX_Log"));
	//		actionAbout_zmX = new QAction(MainWindow);
	//		actionAbout_zmX->setObjectName(QString::fromUtf8("actionAbout_zmX"));
	//		actionOpen_Remote_Support_Session = new QAction(MainWindow);
	//		actionOpen_Remote_Support_Session->setObjectName(QString::fromUtf8("actionOpen_Remote_Support_Session"));
	//		centralwidget = new QWidget(MainWindow);
	//		centralwidget->setObjectName(QString::fromUtf8("centralwidget"));
	//		gridLayout_4 = new QGridLayout(centralwidget);
	//		gridLayout_4->setObjectName(QString::fromUtf8("gridLayout_4"));
	//		mainTabs = new QTabWidget(centralwidget);
	//		mainTabs->setObjectName(QString::fromUtf8("mainTabs"));
	//		tab = new QWidget();
	//		tab->setObjectName(QString::fromUtf8("tab"));
	//		gridLayout = new QGridLayout(tab);
	//		gridLayout->setObjectName(QString::fromUtf8("gridLayout"));
	//		groupBox_3 = new QGroupBox(tab);
	//		groupBox_3->setObjectName(QString::fromUtf8("groupBox_3"));
	//		groupBox_3->setFlat(true);
	//		textBrowser_3 = new QTextBrowser(groupBox_3);
	//		textBrowser_3->setObjectName(QString::fromUtf8("textBrowser_3"));
	//		textBrowser_3->setGeometry(QRect(0, 20, 401, 192));
	//		textBrowser_3->setAutoFillBackground(false);
	//		textBrowser_3->setFrameShape(QFrame::NoFrame);
	//
	//		gridLayout->addWidget(groupBox_3, 1, 0, 1, 1);
	//
	//		groupBox = new QGroupBox(tab);
	//		groupBox->setObjectName(QString::fromUtf8("groupBox"));
	//		groupBox->setFlat(true);
	//		textBrowser = new QTextBrowser(groupBox);
	//		textBrowser->setObjectName(QString::fromUtf8("textBrowser"));
	//		textBrowser->setGeometry(QRect(0, 20, 401, 192));
	//		textBrowser->setAutoFillBackground(false);
	//		textBrowser->setFrameShape(QFrame::NoFrame);
	//
	//		gridLayout->addWidget(groupBox, 0, 0, 1, 1);
	//
	//		groupBox_2 = new QGroupBox(tab);
	//		groupBox_2->setObjectName(QString::fromUtf8("groupBox_2"));
	//		groupBox_2->setFlat(true);
	//		textBrowser_2 = new QTextBrowser(groupBox_2);
	//		textBrowser_2->setObjectName(QString::fromUtf8("textBrowser_2"));
	//		textBrowser_2->setGeometry(QRect(10, 20, 401, 192));
	//		textBrowser_2->setAutoFillBackground(false);
	//		textBrowser_2->setFrameShape(QFrame::NoFrame);
	//
	//		gridLayout->addWidget(groupBox_2, 0, 1, 1, 1);
	//
	//		mainTabs->addTab(tab, QString());
	//		MetadataScanTab = new QWidget();
	//		MetadataScanTab->setObjectName(QString::fromUtf8("MetadataScanTab"));
	//		MetadataScanTab->setEnabled(true);
	//		gridLayout_8 = new QGridLayout(MetadataScanTab);
	//		gridLayout_8->setObjectName(QString::fromUtf8("gridLayout_8"));
	//		label = new QLabel(MetadataScanTab);
	//		label->setObjectName(QString::fromUtf8("label"));
	//		label->setAlignment(Qt::AlignRight|Qt::AlignTrailing|Qt::AlignVCenter);
	//
	//		gridLayout_8->addWidget(label, 0, 2, 1, 1);
	//
	//		checkBox = new QCheckBox(MetadataScanTab);
	//		checkBox->setObjectName(QString::fromUtf8("checkBox"));
	//
	//		gridLayout_8->addWidget(checkBox, 0, 5, 1, 1);
	//
	//		comboBox_2 = new QComboBox(MetadataScanTab);
	//		comboBox_2->addItem(QString());
	//		comboBox_2->addItem(QString());
	//		comboBox_2->setObjectName(QString::fromUtf8("comboBox_2"));
	//		QSizePolicy sizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
	//		sizePolicy.setHorizontalStretch(0);
	//		sizePolicy.setVerticalStretch(0);
	//		sizePolicy.setHeightForWidth(comboBox_2->sizePolicy().hasHeightForWidth());
	//		comboBox_2->setSizePolicy(sizePolicy);
	//
	//		gridLayout_8->addWidget(comboBox_2, 0, 4, 1, 1);
	//
	//		pushButton = new QPushButton(MetadataScanTab);
	//		pushButton->setObjectName(QString::fromUtf8("pushButton"));
	//
	//		gridLayout_8->addWidget(pushButton, 0, 8, 1, 1);
	//
	//		comboBox = new QComboBox(MetadataScanTab);
	//		comboBox->addItem(QString());
	//		comboBox->addItem(QString());
	//		comboBox->setObjectName(QString::fromUtf8("comboBox"));
	//
	//		gridLayout_8->addWidget(comboBox, 0, 3, 1, 1);
	//
	//		tableWidget = new QTableWidget(MetadataScanTab);
	//		if (tableWidget->columnCount() < 5)
	//			tableWidget->setColumnCount(5);
	//		QTableWidgetItem *__qtablewidgetitem = new QTableWidgetItem();
	//		tableWidget->setHorizontalHeaderItem(0, __qtablewidgetitem);
	//		QTableWidgetItem *__qtablewidgetitem1 = new QTableWidgetItem();
	//		tableWidget->setHorizontalHeaderItem(1, __qtablewidgetitem1);
	//		QTableWidgetItem *__qtablewidgetitem2 = new QTableWidgetItem();
	//		tableWidget->setHorizontalHeaderItem(2, __qtablewidgetitem2);
	//		QTableWidgetItem *__qtablewidgetitem3 = new QTableWidgetItem();
	//		tableWidget->setHorizontalHeaderItem(3, __qtablewidgetitem3);
	//		QTableWidgetItem *__qtablewidgetitem4 = new QTableWidgetItem();
	//		tableWidget->setHorizontalHeaderItem(4, __qtablewidgetitem4);
	//		tableWidget->setObjectName(QString::fromUtf8("tableWidget"));
	//
	//		gridLayout_8->addWidget(tableWidget, 1, 0, 2, 9);
	//
	//		mainTabs->addTab(MetadataScanTab, QString());
	//
	//		gridLayout_4->addWidget(mainTabs, 0, 0, 1, 1);
	//
	//		MainWindow->setCentralWidget(centralwidget);
	//		menubar = new QMenuBar(MainWindow);
	//		menubar->setObjectName(QString::fromUtf8("menubar"));
	//		menubar->setGeometry(QRect(0, 0, 1151, 21));
	//		menuFile = new QMenu(menubar);
	//		menuFile->setObjectName(QString::fromUtf8("menuFile"));
	//		menuEdit = new QMenu(menubar);
	//		menuEdit->setObjectName(QString::fromUtf8("menuEdit"));
	//		menuLibrary = new QMenu(menubar);
	//		menuLibrary->setObjectName(QString::fromUtf8("menuLibrary"));
	//		menuSettings = new QMenu(menubar);
	//		menuSettings->setObjectName(QString::fromUtf8("menuSettings"));
	//		MainWindow->setMenuBar(menubar);
	//		statusbar = new QStatusBar(MainWindow);
	//		statusbar->setObjectName(QString::fromUtf8("statusbar"));
	//		MainWindow->setStatusBar(statusbar);
	//		dockWidget_2 = new QDockWidget(MainWindow);
	//		dockWidget_2->setObjectName(QString::fromUtf8("dockWidget_2"));
	//		dockWidget_2->setFloating(false);
	//		dockWidgetContents_2 = new QWidget();
	//		dockWidgetContents_2->setObjectName(QString::fromUtf8("dockWidgetContents_2"));
	//		gridLayout_2 = new QGridLayout(dockWidgetContents_2);
	//		gridLayout_2->setObjectName(QString::fromUtf8("gridLayout_2"));
	//		toolBox = new QToolBox(dockWidgetContents_2);
	//		toolBox->setObjectName(QString::fromUtf8("toolBox"));
	//		page = new QWidget();
	//		page->setObjectName(QString::fromUtf8("page"));
	//		page->setGeometry(QRect(0, 0, 98, 89));
	//		gridLayout_3 = new QGridLayout(page);
	//		gridLayout_3->setObjectName(QString::fromUtf8("gridLayout_3"));
	//		treeWidget = new QTreeWidget(page);
	//		QTreeWidgetItem *__qtreewidgetitem = new QTreeWidgetItem(treeWidget);
	//		new QTreeWidgetItem(__qtreewidgetitem);
	//		QTreeWidgetItem *__qtreewidgetitem1 = new QTreeWidgetItem(treeWidget);
	//		new QTreeWidgetItem(__qtreewidgetitem1);
	//		QTreeWidgetItem *__qtreewidgetitem2 = new QTreeWidgetItem(treeWidget);
	//		new QTreeWidgetItem(__qtreewidgetitem2);
	//		QTreeWidgetItem *__qtreewidgetitem3 = new QTreeWidgetItem(treeWidget);
	//		new QTreeWidgetItem(__qtreewidgetitem3);
	//		treeWidget->setObjectName(QString::fromUtf8("treeWidget"));
	//
	//		gridLayout_3->addWidget(treeWidget, 0, 0, 1, 1);
	//
	//		toolBox->addItem(page, QString::fromUtf8("Servers and Roles"));
	//		page_3 = new QWidget();
	//		page_3->setObjectName(QString::fromUtf8("page_3"));
	//		page_3->setGeometry(QRect(0, 0, 98, 28));
	//		toolBox->addItem(page_3, QString::fromUtf8("Addresses and Domains"));
	//		page_4 = new QWidget();
	//		page_4->setObjectName(QString::fromUtf8("page_4"));
	//		page_4->setGeometry(QRect(0, 0, 98, 28));
	//		toolBox->addItem(page_4, QString::fromUtf8("Configuration and COS"));
	//		page_2 = new QWidget();
	//		page_2->setObjectName(QString::fromUtf8("page_2"));
	//		page_2->setGeometry(QRect(0, 0, 274, 294));
	//		gridLayout_5 = new QGridLayout(page_2);
	//		gridLayout_5->setObjectName(QString::fromUtf8("gridLayout_5"));
	//		reportsDiagTree = new QTreeWidget(page_2);
	//		QTreeWidgetItem *__qtreewidgetitem4 = new QTreeWidgetItem(reportsDiagTree);
	//		new QTreeWidgetItem(__qtreewidgetitem4);
	//		new QTreeWidgetItem(__qtreewidgetitem4);
	//		new QTreeWidgetItem(__qtreewidgetitem4);
	//		reportsDiagTree->setObjectName(QString::fromUtf8("reportsDiagTree"));
	//		reportsDiagTree->setFrameShape(QFrame::NoFrame);
	//		reportsDiagTree->setFrameShadow(QFrame::Plain);
	//
	//		gridLayout_5->addWidget(reportsDiagTree, 0, 0, 1, 1);
	//
	//		toolBox->addItem(page_2, QString::fromUtf8("Reports and Diagnostics"));
	//
	//		gridLayout_2->addWidget(toolBox, 0, 0, 1, 1);
	//
	//		dockWidget_2->setWidget(dockWidgetContents_2);
	//		MainWindow->addDockWidget(Qt::LeftDockWidgetArea, dockWidget_2);
	//		dockWidget_3 = new QDockWidget(MainWindow);
	//		dockWidget_3->setObjectName(QString::fromUtf8("dockWidget_3"));
	//		dockWidget_3->setMinimumSize(QSize(167, 180));
	//		dockWidget_3->setFeatures(QDockWidget::DockWidgetFeatureMask);
	//		dockWidgetContents_4 = new QWidget();
	//		dockWidgetContents_4->setObjectName(QString::fromUtf8("dockWidgetContents_4"));
	//		verticalLayout_2 = new QVBoxLayout(dockWidgetContents_4);
	//		verticalLayout_2->setObjectName(QString::fromUtf8("verticalLayout_2"));
	//		tabWidget_2 = new QTabWidget(dockWidgetContents_4);
	//		tabWidget_2->setObjectName(QString::fromUtf8("tabWidget_2"));
	//		tabWidget_2->setTabPosition(QTabWidget::South);
	//		tabWidget_2->setTabShape(QTabWidget::Rounded);
	//		tab_3 = new QWidget();
	//		tab_3->setObjectName(QString::fromUtf8("tab_3"));
	//		verticalLayout = new QVBoxLayout(tab_3);
	//		verticalLayout->setObjectName(QString::fromUtf8("verticalLayout"));
	//		zmxlogger = new QTextEdit(tab_3);
	//		zmxlogger->setObjectName(QString::fromUtf8("zmxlogger"));
	//		zmxlogger->setFrameShape(QFrame::NoFrame);
	//		zmxlogger->setLineWidth(0);
	//		zmxlogger->setReadOnly(true);
	//
	//		verticalLayout->addWidget(zmxlogger);
	//
	//		tabWidget_2->addTab(tab_3, QString());
	//		tab_4 = new QWidget();
	//		tab_4->setObjectName(QString::fromUtf8("tab_4"));
	//		tabWidget_2->addTab(tab_4, QString());
	//		tab_5 = new QWidget();
	//		tab_5->setObjectName(QString::fromUtf8("tab_5"));
	//		tabWidget_2->addTab(tab_5, QString());
	//
	//		verticalLayout_2->addWidget(tabWidget_2);
	//
	//		dockWidget_3->setWidget(dockWidgetContents_4);
	//		MainWindow->addDockWidget(Qt::BottomDockWidgetArea, dockWidget_3);
	//
	//		menubar->addAction(menuFile->menuAction());
	//		menubar->addAction(menuEdit->menuAction());
	//		menubar->addAction(menuLibrary->menuAction());
	//		menubar->addAction(menuSettings->menuAction());
	//		menuFile->addAction(actionOpen_Settings_Backup);
	//		menuFile->addAction(actionSave_Settings_Backup);
	//		menuFile->addAction(actionSync_Settings);
	//		menuFile->addSeparator();
	//		menuFile->addAction(actionOpen_zmX_Log);
	//		menuFile->addSeparator();
	//		menuEdit->addAction(actionCopy);
	//		menuEdit->addAction(actionCut);
	//		menuEdit->addAction(actionPaste);
	//		menuLibrary->addAction(actionOffline_Wiki);
	//		menuLibrary->addAction(actionZimbra_Support);
	//		menuLibrary->addAction(actionOfficial_Documentation);
	//		menuSettings->addAction(actionZMX_Settings);
	//		menuSettings->addAction(actionOpen_Remote_Support_Session);
	//
	//		retranslateUi(MainWindow);
	//
	//		mainTabs->setCurrentIndex(1);
	//		toolBox->setCurrentIndex(3);
	//		tabWidget_2->setCurrentIndex(0);
	//
	//
	//	QMetaObject::connectSlotsByName(MainWindow);
	//	} // setupUi
	//
	//	void retranslateUi(QMainWindow *MainWindow)
	//	{
	//		MainWindow->setWindowTitle(QCoreApplication::translate("MainWindow", "zmX - CONST_VERSION", nullptr));
	//		#if QT_CONFIG(statustip)
	//		MainWindow->setStatusTip(QString());
	//		#endif // QT_CONFIG(statustip)
	//		actionZMX_Settings->setText(QCoreApplication::translate("MainWindow", "zmX Settings", nullptr));
	//		actionOffline_Wiki->setText(QCoreApplication::translate("MainWindow", "Offline Wiki", nullptr));
	//		actionZimbra_Support->setText(QCoreApplication::translate("MainWindow", "Zimbra Support", nullptr));
	//		actionOfficial_Documentation->setText(QCoreApplication::translate("MainWindow", "Official Documentation", nullptr));
	//		actionCopy->setText(QCoreApplication::translate("MainWindow", "Copy", nullptr));
	//		actionCut->setText(QCoreApplication::translate("MainWindow", "Cut", nullptr));
	//		actionPaste->setText(QCoreApplication::translate("MainWindow", "Paste", nullptr));
	//		actionOpen_Settings_Backup->setText(QCoreApplication::translate("MainWindow", "Open Settings Backup", nullptr));
	//		actionSave_Settings_Backup->setText(QCoreApplication::translate("MainWindow", "Save Settings Backup", nullptr));
	//		actionSync_Settings->setText(QCoreApplication::translate("MainWindow", "Sync Settings", nullptr));
	//		actionOpen_zmX_Log->setText(QCoreApplication::translate("MainWindow", "Open zmX Log", nullptr));
	//		actionAbout_zmX->setText(QCoreApplication::translate("MainWindow", "About zmX", nullptr));
	//		actionOpen_Remote_Support_Session->setText(QCoreApplication::translate("MainWindow", "Open Remote Support Session", nullptr));
	//		groupBox_3->setTitle(QCoreApplication::translate("MainWindow", "Change History", nullptr));
	//		groupBox->setTitle(QCoreApplication::translate("MainWindow", "Latest News", nullptr));
	//		textBrowser->setHtml(QCoreApplication::translate("MainWindow", "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.0//EN\" \"http://www.w3.org/TR/REC-html40/strict.dtd\">\n"
	//		"<html><head><meta name=\"qrichtext\" content=\"1\" /><style type=\"text/css\">\n"
	//		"p, li { white-space: pre-wrap; }\n"
	//		"</style></head><body style=\" font-family:'MS Shell Dlg 2'; font-size:8.25pt; font-weight:400; font-style:normal;\">\n"
	//		"<p style=\" margin-top:0px; margin-bottom:0px; margin-left:0px; margin-right:0px; -qt-block-indent:0; text-indent:0px;\"><span style=\" font-family:'MS Shell Dlg 2';\">TEST</span></p></body></html>", nullptr));
	//		groupBox_2->setTitle(QCoreApplication::translate("MainWindow", "Zimbra Community", nullptr));
	//		mainTabs->setTabText(mainTabs->indexOf(tab), QCoreApplication::translate("MainWindow", "Home", nullptr));
	//		label->setText(QCoreApplication::translate("MainWindow", "Connect via:", nullptr));
	//		#if QT_CONFIG(whatsthis)
	//		checkBox->setWhatsThis(QString());
	//		#endif // QT_CONFIG(whatsthis)
	//		checkBox->setText(QCoreApplication::translate("MainWindow", "Ignore zimbraMailHost", nullptr));
	//		comboBox_2->setItemText(0, QCoreApplication::translate("MainWindow", "...", nullptr));
	//		comboBox_2->setItemText(1, QCoreApplication::translate("MainWindow", "New Item", nullptr));
	//
	//		pushButton->setText(QCoreApplication::translate("MainWindow", "Scan", nullptr));
	//		comboBox->setItemText(0, QCoreApplication::translate("MainWindow", "SSH", nullptr));
	//		comboBox->setItemText(1, QCoreApplication::translate("MainWindow", "HTTPS", nullptr));
	//
	//		QTableWidgetItem *___qtablewidgetitem = tableWidget->horizontalHeaderItem(0);
	//		___qtablewidgetitem->setText(QCoreApplication::translate("MainWindow", "Mailbox Host", nullptr));
	//		QTableWidgetItem *___qtablewidgetitem1 = tableWidget->horizontalHeaderItem(1);
	//		___qtablewidgetitem1->setText(QCoreApplication::translate("MainWindow", "Account", nullptr));
	//		QTableWidgetItem *___qtablewidgetitem2 = tableWidget->horizontalHeaderItem(2);
	//		___qtablewidgetitem2->setText(QCoreApplication::translate("MainWindow", "Group", nullptr));
	//		QTableWidgetItem *___qtablewidgetitem3 = tableWidget->horizontalHeaderItem(3);
	//		___qtablewidgetitem3->setText(QCoreApplication::translate("MainWindow", "Locator", nullptr));
	//		QTableWidgetItem *___qtablewidgetitem4 = tableWidget->horizontalHeaderItem(4);
	//		___qtablewidgetitem4->setText(QCoreApplication::translate("MainWindow", "Metadata", nullptr));
	//		mainTabs->setTabText(mainTabs->indexOf(MetadataScanTab), QCoreApplication::translate("MainWindow", "Page", nullptr));
	//		menuFile->setTitle(QCoreApplication::translate("MainWindow", "File", nullptr));
	//		menuEdit->setTitle(QCoreApplication::translate("MainWindow", "Edit", nullptr));
	//		menuLibrary->setTitle(QCoreApplication::translate("MainWindow", "Library", nullptr));
	//		menuSettings->setTitle(QCoreApplication::translate("MainWindow", "Tools", nullptr));
	//		dockWidget_2->setWindowTitle(QCoreApplication::translate("MainWindow", "Toolbox", nullptr));
	//		QTreeWidgetItem *___qtreewidgetitem = treeWidget->headerItem();
	//		___qtreewidgetitem->setText(0, QCoreApplication::translate("MainWindow", "Server Role", nullptr));
	//
	//		const bool __sortingEnabled = treeWidget->isSortingEnabled();
	//		treeWidget->setSortingEnabled(false);
	//		QTreeWidgetItem *___qtreewidgetitem1 = treeWidget->topLevelItem(0);
	//		___qtreewidgetitem1->setText(0, QCoreApplication::translate("MainWindow", "LDAP", nullptr));
	//		QTreeWidgetItem *___qtreewidgetitem2 = ___qtreewidgetitem1->child(0);
	//		___qtreewidgetitem2->setText(0, QCoreApplication::translate("MainWindow", "LDAP_TEMPALTE", nullptr));
	//		QTreeWidgetItem *___qtreewidgetitem3 = treeWidget->topLevelItem(1);
	//		___qtreewidgetitem3->setText(0, QCoreApplication::translate("MainWindow", "Mailbox", nullptr));
	//		QTreeWidgetItem *___qtreewidgetitem4 = ___qtreewidgetitem3->child(0);
	//		___qtreewidgetitem4->setText(0, QCoreApplication::translate("MainWindow", "MAILBOX_TEMPLATE", nullptr));
	//		QTreeWidgetItem *___qtreewidgetitem5 = treeWidget->topLevelItem(2);
	//		___qtreewidgetitem5->setText(0, QCoreApplication::translate("MainWindow", "MTA", nullptr));
	//		QTreeWidgetItem *___qtreewidgetitem6 = ___qtreewidgetitem5->child(0);
	//		___qtreewidgetitem6->setText(0, QCoreApplication::translate("MainWindow", "MTA_TEMPLATE", nullptr));
	//		QTreeWidgetItem *___qtreewidgetitem7 = treeWidget->topLevelItem(3);
	//		___qtreewidgetitem7->setText(0, QCoreApplication::translate("MainWindow", "Proxy", nullptr));
	//		QTreeWidgetItem *___qtreewidgetitem8 = ___qtreewidgetitem7->child(0);
	//		___qtreewidgetitem8->setText(0, QCoreApplication::translate("MainWindow", "PROXY_TEMPLATE", nullptr));
	//		treeWidget->setSortingEnabled(__sortingEnabled);
	//
	//		toolBox->setItemText(toolBox->indexOf(page), QCoreApplication::translate("MainWindow", "Servers and Roles", nullptr));
	//		toolBox->setItemText(toolBox->indexOf(page_3), QCoreApplication::translate("MainWindow", "Addresses and Domains", nullptr));
	//		toolBox->setItemText(toolBox->indexOf(page_4), QCoreApplication::translate("MainWindow", "Configuration and COS", nullptr));
	//		QTreeWidgetItem *___qtreewidgetitem9 = reportsDiagTree->headerItem();
	//		___qtreewidgetitem9->setText(0, QCoreApplication::translate("MainWindow", "1", nullptr));
	//
	//		const bool __sortingEnabled1 = reportsDiagTree->isSortingEnabled();
	//		reportsDiagTree->setSortingEnabled(false);
	//		QTreeWidgetItem *___qtreewidgetitem10 = reportsDiagTree->topLevelItem(0);
	//		___qtreewidgetitem10->setText(0, QCoreApplication::translate("MainWindow", "Metadata", nullptr));
	//		QTreeWidgetItem *___qtreewidgetitem11 = ___qtreewidgetitem10->child(0);
	//		___qtreewidgetitem11->setText(0, QCoreApplication::translate("MainWindow", "Analyze Metadata", nullptr));
	//		QTreeWidgetItem *___qtreewidgetitem12 = ___qtreewidgetitem10->child(1);
	//		___qtreewidgetitem12->setText(0, QCoreApplication::translate("MainWindow", "Backup Metadata", nullptr));
	//		QTreeWidgetItem *___qtreewidgetitem13 = ___qtreewidgetitem10->child(2);
	//		___qtreewidgetitem13->setText(0, QCoreApplication::translate("MainWindow", "Rebuild Metadata", nullptr));
	//		reportsDiagTree->setSortingEnabled(__sortingEnabled1);
	//
	//		toolBox->setItemText(toolBox->indexOf(page_2), QCoreApplication::translate("MainWindow", "Reports and Diagnostics", nullptr));
	//		dockWidget_3->setWindowTitle(QCoreApplication::translate("MainWindow", "Information", nullptr));
	//		#if QT_CONFIG(accessibility)
	//		tab_3->setAccessibleName(QString());
	//		#endif // QT_CONFIG(accessibility)
	//		tabWidget_2->setTabText(tabWidget_2->indexOf(tab_3), QCoreApplication::translate("MainWindow", "Log", nullptr));
	//		tabWidget_2->setTabText(tabWidget_2->indexOf(tab_4), QCoreApplication::translate("MainWindow", "Performance", nullptr));
	//		tabWidget_2->setTabText(tabWidget_2->indexOf(tab_5), QCoreApplication::translate("MainWindow", "Statistics", nullptr));
	//	} // retranslateUi
	//
	//};
	//
	//namespace Ui {
	//	class MainWindow: public Ui_MainWindow {};
	//} // namespace Ui
	//
	//QT_END_NAMESPACE
	//
	//#endif // ZMX_MAIN_WINDOWHXHQDG_H

}
