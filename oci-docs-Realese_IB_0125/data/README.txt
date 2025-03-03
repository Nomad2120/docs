Чтобы приступить к работе с KalkanCrypt, необходимо выполнить следующие действия:



1) Скопировать libkalkancryptwr-64.so и libkalkancryptwr-64.so.1.1.1 в /usr/lib/. Для этого выполните команду:
	
	sudo cp -f libkalkancryptwr-64.so libkalkancryptwr-64.so.1.1.1 /usr/lib/ 



2) Пройти в папку SDK 2.0\C\Linux\libs_for_linux и скопировать папку kalkancrypt в каталог: /opt/
	
	sudo cp -r kalkancrypt /opt/



3) Затем необходимо добавить в переменную окружения LD_LIBRARY_PATH путь к папке с библиотеками. 
	
	export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/opt/kalkancrypt/:/opt/kalkancrypt/lib/engines


4) Пройдите в папку test и откройте файл test.cpp. Измените пути к ключам: container = ""



5) Запустите MAKEFILE.
 




