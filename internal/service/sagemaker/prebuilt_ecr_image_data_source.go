package sagemaker

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

const (
	// SageMaker Algorithm BlazingText
	sageMakerRepositoryBlazingText = "blazingtext"
	// SageMaker Algorithm DeepAR Forecasting
	sageMakerRepositoryDeepARForecasting = "forecasting-deepar"
	// SageMaker Algorithm Factorization Machines
	sageMakerRepositoryFactorizationMachines = "factorization-machines"
	// SageMaker Algorithm Image Classification
	sageMakerRepositoryImageClassification = "image-classification"
	// SageMaker Algorithm IP Insights
	sageMakerRepositoryIPInsights = "ipinsights"
	// SageMaker Algorithm k-means
	sageMakerRepositoryKMeans = "kmeans"
	// SageMaker Algorithm k-nearest-neighbor
	sageMakerRepositoryKNearestNeighbor = "knn"
	// SageMaker Algorithm Latent Dirichlet Allocation
	sageMakerRepositoryLDA = "lda"
	// SageMaker Algorithm Linear Learner
	sageMakerRepositoryLinearLearner = "linear-learner"
	// SageMaker Algorithm Neural Topic Model
	sageMakerRepositoryNeuralTopicModel = "ntm"
	// SageMaker Algorithm Object2Vec
	sageMakerRepositoryObject2Vec = "object2vec"
	// SageMaker Algorithm Object Detection
	sageMakerRepositoryObjectDetection = "object-detection"
	// SageMaker Algorithm PCA
	sageMakerRepositoryPCA = "pca"
	// SageMaker Algorithm Random Cut Forest
	sageMakerRepositoryRandomCutForest = "randomcutforest"
	// SageMaker Algorithm Semantic Segmentation
	sageMakerRepositorySemanticSegmentation = "semantic-segmentation"
	// SageMaker Algorithm Seq2Seq
	sageMakerRepositorySeq2Seq = "seq2seq"
	// SageMaker Algorithm XGBoost
	sageMakerRepositoryXGBoost = "sagemaker-xgboost"
	// SageMaker Library scikit-learn
	sageMakerRepositoryScikitLearn = "sagemaker-scikit-learn"
	// SageMaker Library Spark ML
	sageMakerRepositorySparkML = "sagemaker-sparkml-serving"
	// SageMaker Library TensorFlow Serving
	sageMakerRepositoryTensorFlowServing = "sagemaker-tensorflow-serving"
	// SageMaker Library TensorFlow Serving EIA
	sageMakerRepositoryTensorFlowServingEIA = "sagemaker-tensorflow-serving-eia"
	// SageMaker Repo MXNet Inference
	sageMakerRepositoryMXNetInference = "mxnet-inference"
	// SageMaker Repo MXNet Inference EIA
	sageMakerRepositoryMXNetInferenceEIA = "mxnet-inference-eia"
	// SageMaker Repo MXNet Training
	sageMakerRepositoryMXNetTraining = "mxnet-training"
	// SageMaker Repo PyTorch Inference
	sageMakerRepositoryPyTorchInference = "pytorch-inference"
	// SageMaker Repo PyTorch Inference EIA
	sageMakerRepositoryPyTorchInferenceEIA = "pytorch-inference-eia"
	// SageMaker Repo PyTorch Training
	sageMakerRepositoryPyTorchTraining = "pytorch-training"
	// SageMaker Repo TensorFlow Inference
	sageMakerRepositoryTensorFlowInference = "tensorflow-inference"
	// SageMaker Repo TensorFlow Inference EIA
	sageMakerRepositoryTensorFlowInferenceEIA = "tensorflow-inference-eia"
	// SageMaker Repo TensorFlow Training
	sageMakerRepositoryTensorFlowTraining = "tensorflow-training"
	// SageMaker Repo HuggingFace TensorFlow Training
	sageMakerRepositoryHuggingFaceTensorFlowTraining = "huggingface-tensorflow-training"
	// SageMaker Repo HuggingFace TensorFlow Inference
	sageMakerRepositoryHuggingFaceTensorFlowInference = "huggingface-tensorflow-inference"
	// SageMaker Repo HuggingFace PyTorch Training
	sageMakerRepositoryHuggingFacePyTorchTraining = "huggingface-pytorch-training"
	// SageMaker Repo HuggingFace PyTorch Inference
	sageMakerRepositoryHuggingFacePyTorchInference = "huggingface-pytorch-inference"
)

// FIXME: get rid of all the maps below and probably of this resource.

// https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html
var sageMakerPrebuiltECRImageIDByRegion_Blazing = map[string]string{}

// https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html
var sageMakerPrebuiltECRImageIDByRegion_DeepAR = map[string]string{}

// https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html
var PrebuiltECRImageIDByRegion_FactorMachines = map[string]string{}

// https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html
var sageMakerPrebuiltECRImageIDByRegion_LDA = map[string]string{}

// https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html
var sageMakerPrebuiltECRImageIDByRegion_XGBoost = map[string]string{}

// https://docs.aws.amazon.com/sagemaker/latest/dg/sagemaker-algo-docker-registry-paths.html
// https://docs.aws.amazon.com/sagemaker/latest/dg/pre-built-docker-containers-scikit-learn-spark.html
var PrebuiltECRImageIDByRegion_SparkML = map[string]string{}

// https://github.com/aws/deep-learning-containers/blob/master/available_images.md
// https://github.com/aws/sagemaker-tensorflow-serving-container
var sageMakerPrebuiltECRImageIDByRegion_DeepLearning = map[string]string{}

// https://github.com/aws/sagemaker-tensorflow-serving-container
var sageMakerPrebuiltECRImageIDByRegion_TensorFlowServing = map[string]string{}

func DataSourcePrebuiltECRImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePrebuiltECRImageRead,
		Schema: map[string]*schema.Schema{
			"repository_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					sageMakerRepositoryBlazingText,
					sageMakerRepositoryDeepARForecasting,
					sageMakerRepositoryFactorizationMachines,
					sageMakerRepositoryImageClassification,
					sageMakerRepositoryIPInsights,
					sageMakerRepositoryKMeans,
					sageMakerRepositoryKNearestNeighbor,
					sageMakerRepositoryLDA,
					sageMakerRepositoryLinearLearner,
					sageMakerRepositoryMXNetInference,
					sageMakerRepositoryMXNetInferenceEIA,
					sageMakerRepositoryMXNetTraining,
					sageMakerRepositoryNeuralTopicModel,
					sageMakerRepositoryObject2Vec,
					sageMakerRepositoryObjectDetection,
					sageMakerRepositoryPCA,
					sageMakerRepositoryPyTorchInference,
					sageMakerRepositoryPyTorchInferenceEIA,
					sageMakerRepositoryPyTorchTraining,
					sageMakerRepositoryRandomCutForest,
					sageMakerRepositoryScikitLearn,
					sageMakerRepositorySemanticSegmentation,
					sageMakerRepositorySeq2Seq,
					sageMakerRepositorySparkML,
					sageMakerRepositoryTensorFlowInference,
					sageMakerRepositoryTensorFlowInferenceEIA,
					sageMakerRepositoryTensorFlowServing,
					sageMakerRepositoryTensorFlowServingEIA,
					sageMakerRepositoryTensorFlowTraining,
					sageMakerRepositoryHuggingFaceTensorFlowTraining,
					sageMakerRepositoryHuggingFaceTensorFlowInference,
					sageMakerRepositoryHuggingFacePyTorchTraining,
					sageMakerRepositoryHuggingFacePyTorchInference,
					sageMakerRepositoryXGBoost,
				}, false),
			},

			"dns_suffix": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"image_tag": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1",
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"registry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"registry_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePrebuiltECRImageRead(d *schema.ResourceData, meta interface{}) error {
	region := meta.(*conns.AWSClient).Region
	if v, ok := d.GetOk("region"); ok {
		region = v.(string)
	}

	suffix := meta.(*conns.AWSClient).DNSSuffix
	if v, ok := d.GetOk("dns_suffix"); ok {
		suffix = v.(string)
	}

	repo := d.Get("repository_name").(string)

	id := ""
	switch repo {
	case sageMakerRepositoryBlazingText,
		sageMakerRepositoryImageClassification,
		sageMakerRepositoryObjectDetection,
		sageMakerRepositorySemanticSegmentation,
		sageMakerRepositorySeq2Seq:
		id = sageMakerPrebuiltECRImageIDByRegion_Blazing[region]
	case sageMakerRepositoryDeepARForecasting:
		id = sageMakerPrebuiltECRImageIDByRegion_DeepAR[region]
	case sageMakerRepositoryLDA:
		id = sageMakerPrebuiltECRImageIDByRegion_LDA[region]
	case sageMakerRepositoryXGBoost:
		id = sageMakerPrebuiltECRImageIDByRegion_XGBoost[region]
	case sageMakerRepositoryScikitLearn, sageMakerRepositorySparkML:
		id = PrebuiltECRImageIDByRegion_SparkML[region]
	case sageMakerRepositoryTensorFlowServing, sageMakerRepositoryTensorFlowServingEIA:
		id = sageMakerPrebuiltECRImageIDByRegion_TensorFlowServing[region]
	case sageMakerRepositoryMXNetInference,
		sageMakerRepositoryMXNetInferenceEIA,
		sageMakerRepositoryMXNetTraining,
		sageMakerRepositoryPyTorchInference,
		sageMakerRepositoryPyTorchInferenceEIA,
		sageMakerRepositoryPyTorchTraining,
		sageMakerRepositoryTensorFlowInference,
		sageMakerRepositoryTensorFlowInferenceEIA,
		sageMakerRepositoryTensorFlowTraining,
		sageMakerRepositoryHuggingFaceTensorFlowTraining,
		sageMakerRepositoryHuggingFaceTensorFlowInference,
		sageMakerRepositoryHuggingFacePyTorchTraining,
		sageMakerRepositoryHuggingFacePyTorchInference:
		id = sageMakerPrebuiltECRImageIDByRegion_DeepLearning[region]
	default:
		id = PrebuiltECRImageIDByRegion_FactorMachines[region]
	}

	if id == "" {
		return fmt.Errorf("no registry ID available for region (%s) and repository (%s)", region, repo)
	}

	d.SetId(id)
	d.Set("registry_id", id)
	d.Set("registry_path", PrebuiltECRImageCreatePath(id, region, suffix, repo, d.Get("image_tag").(string)))
	return nil
}

func PrebuiltECRImageCreatePath(id, region, suffix, repo, imageTag string) string {
	return fmt.Sprintf("%s.dkr.ecr.%s.%s/%s:%s", id, region, suffix, repo, imageTag)
}
