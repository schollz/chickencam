import random
import shutil
import json

from SimpleCV import *
import glob

class Trainer():

    def __init__(self,classes, trainPaths):
        self.classes = classes
        self.trainPaths = trainPaths


    def getExtractors(self):
        hhfe = HueHistogramFeatureExtractor(32)
        ehfe = EdgeHistogramFeatureExtractor(32)
        haarfe = HaarLikeFeatureExtractor(fname='simplecv/haar.txt')
        return [hhfe,ehfe,haarfe]

    def getClassifiers(self,extractors):
        svm = SVMClassifier(extractors)
        tree = TreeClassifier(extractors)
        bayes = NaiveBayesClassifier(extractors)
        knn = KNNClassifier(extractors)
        return [svm,tree,bayes,knn]

    def train(self):
        self.classifiers = self.getClassifiers(self.getExtractors())
        for classifier in self.classifiers:
            classifier.train(self.trainPaths,self.classes,verbose=False)

    def test(self,testPaths):
        for classifier in self.classifiers:
            print classifier.test(testPaths,self.classes,verbose=False)

    def visualizeResults(self,classifier,imgs):
        for img in imgs:
            className = classifier.classify(img)
            img.drawText(className,10,10,fontsize=60,color=Color.BLUE)         
        imgs.show()


def main():
    classes = ['none','chicken']

    print("Moving images...")
    if os.path.exists("images"):
        shutil.rmtree("images")
    os.makedirs("images")
    for c in classes:
        os.makedirs(os.path.join("images",c))
        os.makedirs(os.path.join("images",c,"train"))
        os.makedirs(os.path.join("images",c,"test"))

    for txt in glob.glob("data/*.txt"):
        j = json.load(open(txt,'r'))
        imageName = txt.replace(".txt",".jpg")
        trainOrTest = "train"
        if random.random() < 0.2:
            trainOrTest = "test"
        if 'none' in j['Presence']:
            shutil.copyfile(imageName,os.path.join("images","none",trainOrTest,imageName.split("/")[-1]))
        else:
            shutil.copyfile(imageName,os.path.join("images","chicken",trainOrTest,imageName.split("/")[-1]))

    trainPaths = ['./images/'+c+'/train/' for c in classes ]
    testPaths =  ['./images/'+c+'/test/'   for c in classes ]

    print("Training %d classes..." % len(trainPaths))
    trainer = Trainer(classes,trainPaths)
    trainer.train()
    tree = trainer.classifiers[1]
    
    imgs = ImageSet()
    for p in testPaths:
        imgs += ImageSet(p)
    random.shuffle(imgs)

    print("Testing %d classes with all classifiers..." % len(testPaths))
    trainer.test(testPaths)

    # print(trainer.visualizeResults(tree,imgs))
    tree.save("tree.dat")

    print("Testing with TreeClassifier...")
    classifierFile = 'tree.dat'
    classifier = TreeClassifier.load(classifierFile)
    count = 0
    correct_count = 0
    for path in glob.glob("images/*/test/*jpg"):
        guess = classifier.classify(Image(path))
        print(path,guess)
        if guess in path:
            correct_count += 1
        count += 1
    print("%d correct out of %d = %2.1f accuracy" %(correct_count,count,100.0*correct_count/count))

main()



